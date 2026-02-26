# efizzer

**efizzer** — экспериментальный инструмент для сбора покрытия кода UEFI‑прошивок, инструментированных с помощью Clang Coverage. Проект находится на ранней стадии разработки. В текущей версии поверхностно реализован компонент `efizzer-oracle`, который однократно запускает QEMU со специально подготовленной прошивкой и сохраняет "сырые" данные о покрытии инструментированных модулей.

> ⚠️ **Важно**: На данный момент efizzer позволяет только собрать покрытие для анализа, без организации фаззинг‑цикла.

## Предварительные требования

- Операционная система: Linux 
- Инструменты для сборки `audk` и `qemu`
- `go` версии 1.23.6 и выше

## Внешние зависимости
- [`tidwall/btree`](https://github.com/tidwall/btree.git)

## 1. Сборка QEMU с поддержкой устройства efizzer

Используется форк QEMU ([здесь](https://github.com/stokescat/qemu.git)), в который добавлено специальное устройство `efizzer`, необходимое для сбора покрытия.

```bash
git clone https://github.com/stokescat/qemu.git
cd qemu
mkdir build
cd build
../configure --target-list=x86_64-softmmu
make
```
После сборки в директории `build` появится исполняемый файл `qemu-system-x86_64`. Запомните путь к нему (понадобится позже).

## 2. Инструментирование прошивки

Для сбора покрытия требуется прошивка OVMF, собранная с флагами Clang Coverage. Используется форк репозитория `audk` ([здесь](https://github.com/stokescat/audk.git)) с предварительно инструментированными модулями.

### Клонирование и подготовка

```bash
git clone https://github.com/stokescat/audk.git
cd audk
git submodule update --init --recursive
make -C BaseTools
source edksetup.sh
```

### Настройка флагов компиляции

Отредактируйте файл `Conf/tools_def.txt`. Найдите секцию `CLANGDWARF` и строку:

```
DEBUG_CLANGDWARF_X64_CC_FLAGS         = DEF(CLANGDWARF_ALL_CC_FLAGS) -m64 "-DEFIAPI=__attribute__((ms_abi))" -mno-red-zone -mcmodel=small -fpie -fdirect-access-external-data -Oz -flto DEF(CLANGDWARF_X64_TARGET) -g
```

Замените её на (отключение оптимизаций и LTO, необходимых для корректного покрытия):

```
DEBUG_CLANGDWARF_X64_CC_FLAGS         = DEF(CLANGDWARF_ALL_CC_FLAGS) -m64 "-DEFIAPI=__attribute__((ms_abi))" -mno-red-zone -mcmodel=small -fpie -fdirect-access-external-data -O0 -fno-lto DEF(CLANGDWARF_X64_TARGET) -g
```

### Сборка прошивки

```bash
build -a X64 -t CLANGDWARF -p OvmfPkg/OvmfPkgX64.dsc -b DEBUG
```

После успешной сборки образ прошивки будет находиться по пути:
`Build/OvmfX64/DEBUG_CLANGDWARF/FV/OVMF.fd`

### Инструментированные модули

В текущей версии покрытие собирается для следующих модулей (идентифицируются по GUID):

- `EnhancedFatDxe` (`961578FE-B6B7-44c3-AF35-6BC705CD2B1F`)
- `PciBusDxe` (`93B80004-9FB3-11d4-9A3A-0090273FC14D`)
- `DxeCore` (`D6A2CB7F-6A18-4e2f-B43B-9920A733700A`)
- `Shell` (`7C04A583-9E3E-4f1c-AD65-E05268D0B4D1`)

(?) Возможно, стоит инструментировать всю прошивку сразу или создать инструмент для выбора инструментируемых модулей перед сборкой прошивки (?)


## 3. Настройка путей в efizzer-oracle

Перед сборкой `efizzer-oracle` необходимо указать абсолютные пути к собранному QEMU и прошивке в файле `cmd/efizzer-oracle/main.go`. Найдите и измените следующие переменные ([здесь](https://github.com/stokescat/efizzer/blob/744e46141fc8c0ca971d54bb1b441c83b2fbbdfb/cmd/efizzer-oracle/main.go#L18)):

```go
gWorkDir = ""  // можно оставить пустым – будет использоваться текущая директория
gMachinePath = "/абсолютный/путь/к/qemu/build/qemu-system-x86_64"
firmwarePath = "/абсолютный/путь/к/audk/Build/OvmfX64/DEBUG_CLANGDWARF/FV/OVMF.fd"
```

## 4. Сборка efizzer-oracle

В корне репозитория efizzer выполните:

```bash
make
```

После успешной сборки в директории `bin` появится исполняемый файл `efizzer-oracle`.

## 5. Запуск efizzer-oracle и получение покрытия

Перейдите в папку `bin` и запустите утилиту:

```bash
cd bin
./efizzer-oracle
```

Процесс будет выполняться до тех пор, пока вы не прервёте его нажатием `Ctrl+C`. В течение работы efizzer-oracle однократно запускает QEMU с инструментированной прошивкой, собирает данные о покрытии и сохраняет их в файлы с [расширением `.rawcov`](internal/rawcov/README.md).

### Выходные данные

В директории `bin` появятся файлы следующего вида:

- `<GUID>.rawcov` — например, `961578FE-B6B7-44c3-AF35-6BC705CD2B1F.rawcov`. Содержит «сырые» адреса (без учёта смещения для модулей из FV), соответствующие сработавшим точкам покрытия в конкретном модуле.
- `undefined.rawcov` — адреса, которые не удалось сопоставить ни с одним загруженным модулем. Обычно возникает для модуля `DxeCore`, так как он сообщает о покрытии раньше, чем отправляет событие о своей загрузке. В текущей версии этот файл не обрабатывается повторно и не используется в дальнейшем анализе.

## Примечания о текущем состоянии

- `efizzer-oracle` выполняет **один** запуск QEMU, после чего завершает работу (по Ctrl+C). Никакого автоматического цикла фаззинга или взаимодействия с менеджером пока нет.
- Компоненты `efizzer-manager` и `executor` **не реализованы**; описание в данном README относится только к сборщику покрытия.
- Компонент `efizzer-oracle` также пока реализован поверхностно и требует доработки
