module goTest2

go 1.12

require (
	github.com/Shopify/sarama v1.24.1
	github.com/auth0/go-jwt-middleware v0.0.0-20200810150920-a32d7af194d1 // indirect
	github.com/coreos/etcd v3.3.17+incompatible
	github.com/didip/tollbooth v4.0.2+incompatible
	github.com/didip/tollbooth_gin v0.0.0-20170928041415-5752492be505
	github.com/fatih/pool v3.0.0+incompatible
	github.com/funny/link v0.0.0-20190805113223-98708916287b
	github.com/funny/utest v0.0.0-20161029064919-43870a374500 // indirect
	github.com/gin-gonic/gin v1.4.0
	github.com/go-kit/kit v0.10.0
	github.com/go-redsync/redsync v1.4.2 // indirect
	github.com/go-redsync/redsync/v3 v3.0.0
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/btree v1.0.0 // indirect
	github.com/gorilla/mux v1.7.4
	github.com/gorilla/websocket v1.4.1
	github.com/hashicorp/consul/api v1.5.0
	github.com/itfantasy/gonode v0.0.0-20191022090118-359a2ae7228e
	github.com/juju/ratelimit v1.0.1 // indirect
	github.com/julianshen/gin-limiter v0.0.0-20161123033831-fc39b5e90fe7
	github.com/klauspost/cpuid v1.2.1 // indirect
	github.com/klauspost/reedsolomon v1.9.3 // indirect
	github.com/lzb/replace v0.0.0-00010101000000-000000000000
	github.com/micro/go-micro v1.16.0
	github.com/mkevac/debugcharts v0.0.0-20180124214838-d3203a8fa926 // indirect
	github.com/nareix/joy4 v0.0.0-20181022032202-3ddbc8f9d431
	github.com/opentracing/opentracing-go v1.2.0
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/prometheus/prometheus v1.8.2-0.20200819073411-9438bf735a1e
	github.com/prometheus/tsdb v0.10.0 // indirect
	github.com/pyloque/gocaptain v0.0.0-20160623031443-8e6a4933f710
	github.com/rakyll/statik v0.1.6
	github.com/silenceper/pool v0.0.0-20190419103246-92cc9e6ec7b8
	github.com/smallnest/gofsm v0.0.0-20190306032117-f5ba1bddca7b
	github.com/stretchr/testify v1.5.1
	github.com/stvp/tempredis v0.0.0-20181119212430-b82af8480203
	github.com/templexxx/cpufeat v0.0.0-20180724012125-cef66df7f161 // indirect
	github.com/templexxx/xor v0.0.0-20181023030647-4e92f724b73b // indirect
	github.com/tidwall/gjson v1.6.1
	github.com/tidwall/uhaha v0.2.1
	github.com/tjfoc/gmsm v1.0.1 // indirect
	github.com/unrolled/secure v1.0.4
	github.com/xtaci/kcp-go v5.4.11+incompatible
	github.com/xtaci/lossyconn v0.0.0-20190602105132-8df528c0c9ae // indirect
	go.mongodb.org/mongo-driver v1.3.2
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	gopkg.in/olahol/melody.v1 v1.0.0-20170518105555-d52139073376
	gopkg.in/olivere/elastic.v5 v5.0.82
)

replace github.com/lzb/replace => ./mod-replace/replace

replace google.golang.org/protobuf => google.golang.org/protobuf v1.23.0
