#include <stdlib.h>
#include "rocksdb/c.h"

// This API provides convenient C wrapper functions for rocksdb client.

/* Base */

extern void gorocksdb_destruct_handler(void* state);

/* CompactionFilter */

extern rocksdb_compactionfilter_t* gorocksdb_compactionfilter_create(uintptr_t idx);

/* Comparator */

extern rocksdb_comparator_t* gorocksdb_comparator_create(uintptr_t idx);

extern rocksdb_comparator_t* gorocksdb_comparator_with_ts_create(uintptr_t idx, size_t ts_size);

/* Merge Operator */

extern rocksdb_mergeoperator_t* gorocksdb_mergeoperator_create(uintptr_t idx);
extern void gorocksdb_mergeoperator_delete_value(void* state, const char* v, size_t s);

/* Slice Transform */

extern rocksdb_slicetransform_t* gorocksdb_slicetransform_create(uintptr_t idx);

/* Update Timestamp */

extern size_t get_ts_size_wrapper(void* h, uint32_t cf_id);

extern void gorocksdb_writebatch_update_timestamps(uintptr_t state, rocksdb_writebatch_t* wb, const char* ts, size_t tslen, char** errptr);

extern void gorocksdb_writebatch_wi_update_timestamps(uintptr_t state, rocksdb_writebatch_wi_t* wbwi, const char* ts, size_t tslen, char** errptr);
