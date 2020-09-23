package main

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/tsdb"
	"github.com/prometheus/prometheus/tsdb/tsdbutil"
)

func query2(db *tsdb.DB) {
	querier, err := db.Querier(context.TODO(), 0, time.Now().Unix()+10000)
	if err != nil {
		panic(err)
	}

	seriesSet := query(querier,
		labels.MustNewMatcher(labels.MatchEqual, "guangzhou", "temperature"))
	fmt.Println(seriesSet)
}
func main() {
	var s string
	bytes.NewBufferString(s)
	db := openTestDB(nil, nil)
	query2(db)

	ctx := context.Background()
	app := db.Appender(ctx)

	var ref1, ref2 uint64 = 0, 0
	var err error
	for i := 0; i < 10; i++ {
		// labels的name标识一个度量，类似与mysql中的table. value 标识一个属性或者说向量，类似于mysql中的一个列名，
		// t 标识一个时间戳，这是时序数据库显著的特性，必须带有时间戳，而且以时间戳为唯一主键
		// v 代表一个向量的具体值，类似mysql中一个列名的值. 只能是一个kv的值
		// 这里的例子是，表名是guangzhou，有个温度的属性，然后根据时间戳，分别记录温度在不同时间下的值。这就是时序数据库
		if ref1 != 0 {
			// 快速写入
			err = app.AddFast(ref1, time.Now().Unix()+int64(i), float64(30+i))
			if err != nil {
				panic(err)
			}
		} else {
			// 根据 labels来进行写入数据，会返回这个labels 的引用，有了这个引用就可以快速写入AddFast，
			// 不需要判断是否要新建这个labels
			// 因为时序数据库很多情况下，都是一个labels，但是会对应大量的数据，所以少了总是判断是否需要新建，就会快速一点。
			// 这个与mysql也是一样的，因为mysql有实现定义好了的schema，所以不需要判定表是否存在等各种情况。
			// 但是tsdb不需要用户定义schema，自己会自动判断是否需要新建，所以这个地方可以优化，也就是有 AddFast 这个方法了。
			// prometheus对如何定义和使用这个缓存有很好的实例，参考：scrapeCache
			//
			// 其实prometheus内部已经维护了series cache了，但是很多情况下我们都是在一个goroutine中for循环插入数据，
			// 那么可以直接维护这个ref了，
			// 直接Add只是少了查询series缓存的时间，其实也不是很费时，但是如果能够明确获取ref还是可以加速一些性能的。
			ref1, err = app.Add(labels.FromStrings("guangzhou", "temperature"),
				time.Now().Unix()+int64(i), float64(30+i))
			if err != nil {
				panic(err)
			}
		}

		if ref2 != 0 {
			err = app.AddFast(ref2, time.Now().Unix()+int64(i), float64(30+i))
			if err != nil {
				panic(err)
			}
		} else {
			ref2, err = app.Add(labels.FromStrings("guangzhou", "wind"),
				time.Now().Unix()+int64(i), float64(30+i))
			if err != nil {
				panic(err)
			}
		}

	}

	err = app.Commit()
	if err != nil {
		panic(err)
	}
	// 删除
	query2(db)
	db.Delete(0, time.Now().Unix()+9999,
		labels.MustNewMatcher(labels.MatchEqual, "guangzhou", "wind"))

	querier, err := db.Querier(context.TODO(),
		time.Now().Unix()+5, time.Now().Unix()+100)

	seriesSet := query(querier,
		labels.MustNewMatcher(labels.MatchEqual, "guangzhou", "temperature"))
	fmt.Println(seriesSet)
}

func openTestDB(opts *tsdb.Options, rngs []int64) (db *tsdb.DB) {
	dir := "./tsdb/data"
	//if err := os.MkdirAll("dir", 0777); err != nil {
	//	panic(err)
	//}
	////tmpdir, err := ioutil.TempDir("./", "test")
	////if err != nil {
	////	panic(err)
	////}
	if len(rngs) == 0 {
		db, err := tsdb.Open(dir, nil, nil, opts)
		if err != nil {
			panic(err)
		}
		return db
	}

	return db
}

// 可以使用 tsdbutil.Sample
type sample struct {
	t int64
	v float64
}

func newSample(t int64, v float64) tsdbutil.Sample { return sample{t, v} }
func (s sample) T() int64                          { return s.t }
func (s sample) V() float64                        { return s.v }

// query runs a matcher query against the querier and fully expands its data.
func query(q storage.Querier, matchers ...*labels.Matcher) map[string][]tsdbutil.Sample {
	ss := q.Select(false, nil, matchers...)
	defer func() {
		q.Close()
	}()

	result := map[string][]tsdbutil.Sample{}
	for ss.Next() {
		series := ss.At()

		samples := []tsdbutil.Sample{}
		it := series.Iterator()
		for it.Next() {
			t, v := it.At()
			samples = append(samples, sample{t: t, v: v})
		}

		if len(samples) == 0 {
			continue
		}

		name := series.Labels().String()
		result[name] = samples
	}

	return result
}
