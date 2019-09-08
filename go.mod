module gearbox

require (
	github.com/Azure/azure-sdk-for-go v30.1.0+incompatible
	github.com/Azure/go-autorest/autorest v0.2.0
	github.com/Azure/go-autorest/autorest/to v0.2.0
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/apcera/util v0.0.0-20180322191801-7a50bc84ee48
	github.com/aws/aws-sdk-go v1.20.2
	github.com/cavaliercoder/grab v2.0.0+incompatible
	github.com/cenkalti/backoff v2.1.1+incompatible // indirect

	github.com/clbanning/checkjson v0.0.0-20190418161636-abd3ee163e3e
	github.com/cznic/golex v0.0.0-20181122101858-9c343928389c // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/digitalocean/godo v1.16.0
	github.com/docker/docker v1.13.1 // indirect
	github.com/docker/machine v0.16.1
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/exoscale/egoscale v0.18.1
	github.com/fatih/color v1.7.0
	github.com/fhmq/hmq v0.0.0-20190424074534-daf4a0e0f564
	github.com/gearboxworks/go-jsoncache v1.0.0
	github.com/gearboxworks/go-osbridge v0.0.0-20190605062119-0e1c68c1c70f
	github.com/gearboxworks/go-status v0.0.0-20190528175348-42860fb9e78f
	github.com/gearboxworks/go-systray v0.0.0-20190628045254-f866182abfa7
	github.com/gedex/inflector v0.0.0-20170307190818-16278e9db813
	github.com/gernest/wow v0.1.0
	github.com/getlantern/errors v0.0.0-20190325191628-abdb3e3e36f7
	github.com/getlantern/golog v0.0.0-20190830074920-4ef2e798c2d7 // indirect
	github.com/go-bindata/go-bindata v1.0.0
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/uuid v1.1.1
	github.com/gotk3/gotk3 v0.0.0-20190614104930-c157952b53bd // indirect
	github.com/grandcat/zeroconf v0.0.0-20190424104450-85eadb44205c
	github.com/hashicorp/go-version v1.1.0
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/intel-go/cpuid v0.0.0-20181003105527-1a4a6f06a1c6
	github.com/jinzhu/copier v0.0.0-20180308034124-7e38e58719c3
	github.com/kardianos/service v1.0.0
	github.com/kisielk/gotool v1.0.0 // indirect
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.2.8
	github.com/mattn/go-colorable v0.1.0 // indirect
	github.com/mattn/go-isatty v0.0.4 // indirect
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/miekg/dns v1.1.13 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	github.com/nsf/jsondiff v0.0.0-20190302080047-dbf513526b7f
	github.com/olebedev/emitter v0.0.0-20190110104742-e8d1457e6aee
	github.com/racker/perigee v0.1.0 // indirect
	github.com/rackspace/gophercloud v1.0.0
	github.com/radovskyb/watcher v1.0.6
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/sebest/logrusly v0.0.0-20180315190218-3235eccb8edc
	github.com/segmentio/go-loggly v0.5.0 // indirect
	github.com/shirou/gopsutil v2.18.12+incompatible
	github.com/shirou/w32 v0.0.0-20160930032740-bb4de0191aa4 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3 // indirect
	github.com/sqweek/dialog v0.0.0-20190609154315-3cf53be95497
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v0.0.0-20170224212429-dcecefd839c4 // indirect
	github.com/vmware/govcloudair v0.0.2
	github.com/vmware/govmomi v0.20.1
	github.com/z7zmey/php-parser v0.6.0
	github.com/zserge/lorca v0.1.7
	github.com/zserge/webview v0.0.0-20190123072648-16c93bcaeaeb
	golang.org/x/crypto v0.0.0-20190422183909-d864b10871cd
	golang.org/x/net v0.0.0-20190503192946-f4e77d36d62c
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	google.golang.org/api v0.6.0
	google.golang.org/grpc v1.20.1
	gopkg.in/cheggaaa/pb.v1 v1.0.28
	gopkg.in/guregu/null.v2 v2.1.2 // indirect
)

replace github.com/getlantern/systray => github.com/gearboxworks/go-systray v0.0.0-20190626020534-3518af45bf7c
