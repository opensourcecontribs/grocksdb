package grocksdb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteBatchWithTS(t *testing.T) {
	db, cfs, cleanup := newTestDBMultiCF(t, []string{"default"}, func(opts *Options) {
		opts.SetComparator(newDefaultComparatorWithTS())
	})
	defer cleanup()

	defaultCF := cfs[0]

	var (
		givenKey1 = []byte("key1")
		givenVal1 = []byte("val1")
		givenKey2 = []byte("key2")

		givenTs1 = marshalTimestamp(1)
		givenTs2 = marshalTimestamp(2)
	)
	wo := NewDefaultWriteOptions()
	require.Nil(t, db.PutWithTS(wo, givenKey2, givenTs1, []byte("foo")))

	// create and fill the write batch
	wb := NewWriteBatch()
	defer wb.Destroy()
	wb.PutCFWithTS(defaultCF, givenKey1, givenTs2, givenVal1)
	wb.DeleteCFWithTS(defaultCF, givenKey2, givenTs2)
	require.EqualValues(t, wb.Count(), 2)

	// perform the batch
	require.Nil(t, db.Write(wo, wb))

	// check changes
	ro := NewDefaultReadOptions()
	ro.SetTimestamp(givenTs2)
	v1, t1, err := db.GetWithTS(ro, givenKey1)
	defer v1.Free()
	require.Nil(t, err)
	require.EqualValues(t, v1.Data(), givenVal1)
	require.EqualValues(t, t1.Data(), givenTs2)

	v2, t2, err := db.GetWithTS(ro, givenKey2)
	defer v2.Free()
	require.Nil(t, err)
	require.True(t, v2.Data() == nil)
	require.True(t, t2.Data() == nil)

	wb.Clear()
	// DeleteRange not supported for timestamp
}

func TestWriteBatchIteratorWithTS(t *testing.T) {
	_, cfs, cleanup := newTestDBMultiCF(t, []string{"default"}, func(opts *Options) {
		opts.SetComparator(newDefaultComparatorWithTS())
	})
	defer cleanup()

	defaultCF := cfs[0]

	var (
		givenKey1 = []byte("key1")
		givenVal1 = []byte("val1")
		givenKey2 = []byte("key2")

		givenTs1 = marshalTimestamp(1)
		givenTs2 = marshalTimestamp(2)

		expectedKeyWithTS1 = append(givenKey1, givenTs1...)
		expectedKeyWithTS2 = append(givenKey2, givenTs2...)
	)

	// create and fill the write batch
	wb := NewWriteBatch()
	defer wb.Destroy()
	wb.PutCFWithTS(defaultCF, givenKey1, givenTs1, givenVal1)
	wb.DeleteCFWithTS(defaultCF, givenKey2, givenTs2)
	require.EqualValues(t, wb.Count(), 2)

	// iterate over the batch
	iter := wb.NewIterator()
	require.True(t, iter.Next())
	record := iter.Record()
	require.EqualValues(t, record.Type, WriteBatchValueRecord)
	require.EqualValues(t, record.Key, expectedKeyWithTS1)
	require.EqualValues(t, record.Value, givenVal1)

	require.True(t, iter.Next())
	record = iter.Record()
	require.EqualValues(t, record.Type, WriteBatchDeletionRecord)
	require.EqualValues(t, record.Key, expectedKeyWithTS2)

	// there shouldn't be any left
	require.False(t, iter.Next())
}

func TestWriteBatchWithTSDelayed(t *testing.T) {
	db, cfs, cleanup := newTestDBMultiCF(t, []string{"default", "other"}, func(opts *Options) {
		opts.SetComparator(newDefaultComparatorWithTS())
	})
	defer cleanup()

	defaultCF := cfs[0]
	otherCF := cfs[1]

	_ = defaultCF
	_ = otherCF
	_ = db

	var (
		givenKey1 = []byte("key1")
		givenVal1 = []byte("val1")

		// givenTs1 = marshalTimestamp(1)
		givenTs2 = marshalTimestamp(2)
	)

	var err error

	wb := NewWriteBatchWIWithParams(newDefaultComparatorWithTS(), 0, true, 0, 0)
	defer wb.Destroy()
	err = wb.Put(givenKey1, givenVal1)
	require.Nil(t, err)

	err = wb.UpdateTimestamps(givenTs2, func(cfID uint32) uint64 {
		return 8
	})
	require.Nil(t, err)

	// ro := NewDefaultReadOptions()
	// ro.SetTimestamp(givenTs2)

	// i := db.NewIteratorCF(ro, otherCF)
	// defer i.Close()

	// for i.SeekToFirst(); i.Valid(); i.Next() {
	// 	fmt.Println("=-> key", string(i.Key().Data()))
	// 	fmt.Println("=-> value", string(i.Value().Data()))
	// }

	// fmt.Println("errr", i.Err())

	wo := NewDefaultWriteOptions()
	// perform the batch
	require.Nil(t, db.WriteWI(wo, wb))

	// ro = NewDefaultReadOptions()
	// ro.SetTimestamp(givenTs2)

	// fmt.Println("======>")
	// // check changes
	// ro = NewDefaultReadOptions()
	// ro.SetTimestamp(givenTs1)
	// v1, t1, err := db.GetCFWithTS(ro, otherCF, givenKey1)
	// defer v1.Free()
	// require.Nil(t, err)
	// fmt.Println("=-> value", string(v1.Data()))
	// fmt.Println("=-> timestamp", unmarshalTimestamp(t1.Data()))
	// fmt.Println("======>")

	// // check changes
	// ro = NewDefaultReadOptions()
	// ro.SetTimestamp(givenTs2)
	// v1, t1, err = db.GetCFWithTS(ro, otherCF, givenKey1)
	// defer v1.Free()
	// require.Nil(t, err)
	// fmt.Println("=-> value", string(v1.Data()))
	// fmt.Println("=-> timestamp", unmarshalTimestamp(t1.Data()))
	// fmt.Println("======>")

	// // v2, t2, err := db.GetCFWithTS(ro, otherCF, givenKey2)
	// // defer v2.Free()
	// // require.Nil(t, err)
	// // require.True(t, v2.Data() == nil)
	// // require.True(t, t2.Data() == nil)

	// fmt.Println("=------------------")

	ro := NewDefaultReadOptions()
	ro.SetTimestamp(givenTs2)
	i := db.NewIteratorCF(ro, defaultCF)
	defer i.Close()

	for i.SeekToFirst(); i.Valid(); i.Next() {
		fmt.Println("=-> key", string(i.Key().Data()))
		fmt.Println("=-> value", string(i.Value().Data()))
		fmt.Println("=-> timestamp", unmarshalTimestamp(i.Timestamp().Data()))
	}

	fmt.Println("errr", i.Err())

	wb.Clear()
	// DeleteRange not supported for timestamp
}
