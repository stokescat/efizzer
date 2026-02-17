# rawcov
`rawcov` is a Go package for reading, writing, and merging files in the **.rawcov** format – a simple binary format designed to store code coverage data.

---

## Format Specification

A `.rawcov` file consists of a 16‑byte header followed by an array of fixed‑size records. All multi‑byte integers are stored in **little‑endian** order.

### Header (16 bytes)

| Offset | Size | Field      | Description |
|--------|------|------------|-------------|
| 0      | 6    | Signature  | ASCII characters `'R'`, `'A'`, `'W'`, `'C'`, `'O'`, `'V'` |
| 6      | 2    | Flags      | Reserved for future use; **must be zero** in the current version. |
| 8      | 8    | Count      | Number of records in the file (size of the record array). |

### Record (12 bytes)

| Offset | Size | Field | Description                                                               |
| ------ | ---- | ----- | ------------------------------------------------------------------------- |
| 0      | 8    | Value | An identifier – typically an offset in the module or an absolute address. |
| 8      | 4    | Hit   | Number of times this location was hit (saturates at `0xFFFFFFFF`).        |

**Important:** Records in the file must be **sorted in ascending order** by the `Value` field. This invariant enables efficient merging of multiple coverage files.

---

## Package API

### Types

```go
type RawcovRecord struct {
    Value uint64
    Hit   uint32
}
```
Represents a single coverage record.

```go
type RawcovFile struct {
    // contains unexported fields
}
```
The main type for interacting with a `.rawcov` file. It implements the `RawcovReader` interface.

```go
type RawcovReader interface {
    IsEmpty() bool                  // returns true if the reader has no records
    Get()     (RawcovRecord, error) // returns the next record or ErrNoRecord
    Reset()   error                 // resets the reader to the beginning
    Len()     uint64                 // returns the total number of records
}
```

### Functions

#### `Open(path string) (*RawcovFile, error)`

Opens an existing `.rawcov` file or creates a new empty object if the file does not exist.  
If the file exists but has an invalid signature or non‑zero flags, an appropriate error (`ErrInvSignature`, `ErrInvFlags`) is returned.

### Methods

#### `(*RawcovFile) Close() error`

Closes the associated file (if any) and resets the internal state. After calling `Close`, the object becomes empty and can't be reused. It is safe to call `Close` multiple times.

#### `(*RawcovFile) Get() (RawcovRecord, error)`

Returns the next record from the file. When no more records are available, it returns `ErrNoRecord`.

#### `(*RawcovFile) Reset() error`

Resets the read position to the beginning of the file, allowing records to be read again. If the underlying file is closed, an error may occur.

#### `(*RawcovFile) Len() uint64`

Returns the total number of records in the file.

#### `(*RawcovFile) IsEmpty() bool`

Returns `true` if the file contains zero records.

#### `(*RawcovFile) Merge(src RawcovReader, transform func(RawcovRecord) RawcovRecord) (uint64, error)`

Merges the records from `src` into the current file (`dst`). The merge respects the sorted order: records from both sources are combined, and when the same `Value` appears in both, their `Hit` counts are summed (saturated at `0xFFFFFFFF`).  
The optional `transform` function can modify each record read from `src` before merging. If `transform` is `nil`, the identity function is used.

The merge operation  it creates a temporary file, writes the merged data, and then replaces the original file with the new one.

Returns the number of **new unique values** that were added from `src` (i.e., values that did not exist in the destination before the merge).

### Errors

The package defines several sentinel errors that can be returned:

- `ErrNoRecord` – no more records available (end of file).
- `ErrNil` – method called on a nil `RawcovFile` receiver.
- `ErrArgs` – invalid arguments (e.g., `nil` source or destination equals source).
- `ErrInvSignature` – file signature is not `"RAWCOV"`.
- `ErrInvFlags` – reserved flags field is non‑zero.
- `ErrNoFile` – operation requires an open file but none is associated.

---

## Notes

- Always call `Close` on a `RawcovFile` when you are finished with it to release the underlying file descriptor. The object can still be used after closing (it will behave like an empty file), but file operations will fail.
- All methods except `Merge` are safe for concurrent **read** access only. Concurrent calls to `Get` or `Reset` from multiple goroutines are not synchronized. For concurrent use, external locking is required.
- The `Merge` operation is **not** safe to call concurrently on the same `RawcovFile` instance.