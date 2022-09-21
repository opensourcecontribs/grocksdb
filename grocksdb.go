package grocksdb

// #include <stdint.h>
import "C"
import "runtime/cgo"

//export gorocksdb_get_ts_size
func gorocksdb_get_ts_size(h C.uintptr_t, p1 C.uint32_t) C.size_t {
	fn := cgo.Handle(h).Value().(func(uint32) uint64)
	return C.size_t(fn(uint32(p1)))
}
