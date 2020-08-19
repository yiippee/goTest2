package main

import (
	"context"
	"fmt"
	"github.com/prometheus/prometheus/storage"
	"io/ioutil"
	"time"

	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/tsdb"
	"github.com/prometheus/prometheus/tsdb/tsdbutil"
)

func main() {
	db := openTestDB(nil, nil)
	ctx := context.Background()
	app := db.Appender(ctx)

	for i := 0; i < 10; i++ {
		// labels的name标识一个度量，类似与mysql中的table. value 标识一个属性或者说向量，类似于mysql中的一个列名，
		// t 标识一个时间戳，这是时序数据库显著的特性，必须带有时间戳，而且以时间戳为唯一主键
		// v 代表一个向量的具体值，类似mysql中一个列名的值
		// 这里的例子是，表名是guangzhou，有个温度的属性，然后根据时间戳，分别记录温度在不同时间下的值。这就是时序数据库
		_, err := app.Add(labels.FromStrings("guangzhou", "temperature"), time.Now().Unix()+int64(i), float64(30+i))
		if err != nil {
			panic(err)
		}
	}

	err := app.Commit()
	if err != nil {
		panic(err)
	}

	querier, err := db.Querier(context.TODO(), time.Now().Unix()+5, time.Now().Unix()+100)

	seriesSet := query(querier, labels.MustNewMatcher(labels.MatchEqual, "guangzhou", "temperature"))
	fmt.Println(seriesSet)
}

func openTestDB(opts *tsdb.Options, rngs []int64) (db *tsdb.DB) {
	tmpdir, err := ioutil.TempDir("./", "test")
	if err != nil {
		panic(err)
	}
	if len(rngs) == 0 {
		db, err = tsdb.Open(tmpdir, nil, nil, opts)
	}

	return db
}

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
