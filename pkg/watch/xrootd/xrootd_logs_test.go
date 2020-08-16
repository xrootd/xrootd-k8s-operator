package xrootd

import (
	"bytes"
	"io"
	"regexp"
	"testing"

	"github.com/msoap/byline"
	"github.com/pkg/errors"
)

func testGrepReader(t *testing.T, pattern string, logs string) {
	reader := bytes.NewBufferString(logs)
	regex := regexp.MustCompile(pattern)
	lineReader := byline.NewReader(reader)
	buffer := make([]byte, 50)
	t.Log("Reading...", "regex", regex)
	read, err := lineReader.GrepByRegexp(regex).Read(buffer)
	if err != nil {
		if err == io.EOF {
			err = nil
		} else {
			err = errors.Wrap(err, "unable to read")
			return
		}
	}
	if err != nil {
		t.Error(err)
	}
	t.Logf("read: %d, buffer: %s", read, string(buffer))
	if read == 0 {
		t.Error("Failed pattern match!")
	}
}

func TestGrepReaderForWorkerComponent(t *testing.T) {
	testGrepReader(t, logPatternXrootdWorkerIsConnected, `
+ service=
+ getopts hs: c
+ case $c in
+ service=cmsd
+ getopts hs: c
+ '[' 2 -eq 0 ']'
+ '[' -z cmsd ']'
+ printf '%s\n' cmsd
+ egrep -qv '^(xrootd|cmsd)$'
+ shift 2
+ CONFIG_DIR=/config-etc
+ XROOTD_CONFIG=/config-etc/xrootd.cf
+ XROOTD_RDR_DN=base-xrootd-xrootd-redirector
+ hostname
+ egrep '^base-xrootd-xrootd-redirector-[0-9]+$'
+ COMPONENT_NAME=worker
+ export COMPONENT_NAME
++ whoami
+ echo 'Start cmsd as xrootd user'
+ cmd='cmsd -c /config-etc/xrootd.cf -n worker'
+ exec cmsd -c /config-etc/xrootd.cf -n worker
Start cmsd as xrootd user
200816 12:46:28 001 Starting on Linux 5.7.11-arch1-1
Copr.  2004-2012 Stanford University, xrd version v4.11.2
Config warning: this hostname, base-xrootd-xrootd-worker-0, is registered without a domain qualification.
++++++ cmsd worker@base-xrootd-xrootd-worker-0 initialization started.
Config using configuration file /config-etc/xrootd.cf
=====> xrd.sched mint 32 maxt 8192 avlt 512 idle 780
=====> xrd.network buffsz 0 nodnr
=====> all.adminpath /tmp/xrd
=====> xrd.timeout idle 48h
=====> xrd.network dyndns
=====> xrd.network cache 0
=====> xrd.port 1094
=====> xrd.port 2131
Config maximum number of connections restricted to 1048576
Copr.  2007 Stanford University/SLAC cmsd.
++++++ worker@base-xrootd-xrootd-worker-0 phase 1 initialization started.
=====> all.role server
=====> cms.space 1k 2k
=====> cms.sched cpu 10 io 10 space 80
=====> oss.defaults nomig nodread stage nocheck norcreate
=====> all.export /data r/w nocheck norcreate
=====> all.adminpath /tmp/xrd
=====> all.manager base-xrootd-xrootd-redirector:2131
The following paths are available to the redirector:
ws /data

------ worker@base-xrootd-xrootd-worker-0 phase 1 server initialization completed.
++++++ worker@base-xrootd-xrootd-worker-0 phase 2 server initialization started.
Config warning: adminpath resides in /tmp and may be unstable!
++++++ Storage system initialization started.
=====> oss.defaults nomig nodread stage nocheck norcreate
=====> oss.alloc 512M 2 0
=====> oss.fdlimit * max
=====> all.export /data r/w nocheck norcreate
++++++ Configuring standalone mode . . .
Config effective /config-etc/xrootd.cf oss configuration:
       oss.alloc        536870912 2 0
       oss.cachescan    600
       oss.fdlimit      524288 1048576
       oss.maxsize      0
       oss.trace        0
       oss.xfr          1 deny 10800 keep 1200
       oss.memfile off  max 8230289408
       oss.defaults  r/w  nocheck nodread nomig norcreate nopurge stage xattr
       oss.path /data r/w  nocheck nodread nomig norcreate nopurge stage xattr
------ Storage system initialization completed.
200816 12:46:28 001 Meter: Found 1 filesystem(s); 107GB total (87% util); 14GB free (14GB max)
------ worker@base-xrootd-xrootd-worker-0 phase 2 server initialization completed.
------ cmsd worker@base-xrootd-xrootd-worker-0:40701 initialization completed.
200816 12:46:28 033 Start: Waiting for primary server to login.
200816 12:46:28 035 do_Login:: Primary server 1 logged in; data port is 1094
Config Connecting to 1 manager and 1 site.
200816 12:46:28 015 Config: server service enabled.
200816 12:46:28 036 State: Status changed to active + staging
200816 12:46:30 017 XrdOpen: Unable to create socket for 'base-xrootd-xrootd-redirector'; dynamic name or service not yet registered
200816 12:46:54 017 ManTree: Now connected to 1 root node(s)
200816 12:46:54 017 Protocol: Logged into 192.168.82.63
`,
	)
}

func TestGrepReaderForRedirectorComponent(t *testing.T) {
	testGrepReader(t, logPatternXrootdRedirectorIsConnected, `
+ service=
+ getopts hs: c
+ case $c in
+ service=cmsd
+ getopts hs: c
+ '[' 2 -eq 0 ']'
+ '[' -z cmsd ']'
+ printf '%s\n' cmsd
+ egrep -qv '^(xrootd|cmsd)$'
+ shift 2
+ CONFIG_DIR=/config-etc
+ XROOTD_CONFIG=/config-etc/xrootd.cf
+ XROOTD_RDR_DN=base-xrootd-xrootd-redirector
+ hostname
+ egrep '^base-xrootd-xrootd-redirector-[0-9]+$'
base-xrootd-xrootd-redirector-0
+ COMPONENT_NAME=manager
+ export COMPONENT_NAME
++ whoami
Start cmsd as xrootd user
+ echo 'Start cmsd as xrootd user'
+ cmd='cmsd -c /config-etc/xrootd.cf -n manager'
+ exec cmsd -c /config-etc/xrootd.cf -n manager
200816 12:46:28 001 Starting on Linux 5.7.11-arch1-1
Copr.  2004-2012 Stanford University, xrd version v4.11.2
Config warning: this hostname, base-xrootd-xrootd-redirector-0, is registered without a domain qualification.
++++++ cmsd manager@base-xrootd-xrootd-redirector-0 initialization started.
Config using configuration file /config-etc/xrootd.cf
=====> all.adminpath /tmp/xrd
=====> xrd.timeout idle 48h
=====> xrd.network dyndns
=====> xrd.network cache 0
=====> xrd.port 1094
=====> xrd.port 2131
Config maximum number of connections restricted to 1048576
Copr.  2007 Stanford University/SLAC cmsd.
++++++ manager@base-xrootd-xrootd-redirector-0 phase 1 initialization started.
=====> all.role manager
=====> cms.delay servers 1 startup 10 lookup 1 qdl 1
=====> cms.sched cpu 10 io 10 space 80
=====> all.export /data r/w nocheck norcreate
=====> all.adminpath /tmp/xrd
=====> all.manager base-xrootd-xrootd-redirector:2131
------ manager@base-xrootd-xrootd-redirector-0 phase 1 manager initialization completed.
++++++ manager@base-xrootd-xrootd-redirector-0 phase 2 manager initialization started.
Config warning: adminpath resides in /tmp and may be unstable!
++++++ Storage system initialization started.
=====> all.export /data r/w nocheck norcreate
++++++ Configuring standalone mode . . .
Config effective /config-etc/xrootd.cf oss configuration:
       oss.alloc        0 0 0
       oss.cachescan    600
       oss.fdlimit      524288 1048576
       oss.maxsize      0
       oss.trace        0
       oss.xfr          1 deny 10800 keep 1200
       oss.memfile off  max 8230289408
       oss.defaults  r/w  nocheck nodread nomig norcreate nopurge nostage xattr
       oss.path /data r/w  nocheck nodread nomig norcreate nopurge nostage xattr
------ Storage system initialization completed.
------ manager@base-xrootd-xrootd-redirector-0 phase 2 manager initialization completed.
------ cmsd manager@base-xrootd-xrootd-redirector-0:2131 initialization completed.
200816 12:46:36 017 Protocol: redirector.1:19@base-xrootd-xrootd-redirector-0 logged in.
200816 12:46:38 016 Config: manager service enabled.
200816 12:46:38 031 State: Status changed to suspended + nostaging
200816 12:46:54 031 State: Status changed to active + staging
200816 12:46:54 032 Protocol: Primary server.1:20@192-168-82-48.base-xrootd-xrootd-worker.default.svc.cluster.local:1094 logged in.
200816 12:46:54 032 Protocol: server.1:20@192-168-82-48.base-xrootd-xrootd-worker.default.svc.cluster.local:1094 system ID: s-worker@base-xrootd-xrootd-worker-0 2131base-xrootd-xrootd-redirector
=====> Routing for 192.168.82.48: local pub4 prv4
=====> Route all4: 192.168.82.48 Dest=[::192.168.82.48]:1094
	`,
	)
}
