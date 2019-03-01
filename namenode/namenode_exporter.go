package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
)

const (
	namespace = "namenode"
)

var (
	listenAddress  = flag.String("web.listen-address", ":9070", "Address on which to expose metrics and web interface.")
	metricsPath    = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	namenodeJmxUrl = flag.String("namenode.jmx.url", "http://localhost:50070/jmx", "Hadoop JMX URL.")
)

type Exporter struct {
	url                      string
	MissingBlocks            prometheus.Gauge
	UnderReplicatedBlocks    prometheus.Gauge
	CapacityTotal            prometheus.Gauge
	CapacityUsed             prometheus.Gauge
	CapacityRemaining        prometheus.Gauge
	CapacityUsedNonDFS       prometheus.Gauge
	BlocksTotal              prometheus.Gauge
	FilesTotal               prometheus.Gauge
	CorruptBlocks            prometheus.Gauge
	ExcessBlocks             prometheus.Gauge
	StaleDataNodes           prometheus.Gauge
	pnGcCount                prometheus.Gauge
	pnGcTime                 prometheus.Gauge
	cmsGcCount               prometheus.Gauge
	cmsGcTime                prometheus.Gauge
	heapMemoryUsageCommitted prometheus.Gauge
	heapMemoryUsageInit      prometheus.Gauge
	heapMemoryUsageMax       prometheus.Gauge
	heapMemoryUsageUsed      prometheus.Gauge
	lastHATransitionTime     prometheus.Gauge
	state                    prometheus.Gauge
	isActive                 prometheus.Gauge
}

func NewExporter(url string) *Exporter {
	return &Exporter{
		url: url,
		MissingBlocks: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "MissingBlocks",
			Help:      "MissingBlocks",
		}),
		UnderReplicatedBlocks: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "UnderReplicatedBlocks",
			Help:      "UnderReplicatedBlocks",
		}),
		CapacityTotal: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "CapacityTotal",
			Help:      "CapacityTotal",
		}),
		CapacityUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "CapacityUsed",
			Help:      "CapacityUsed",
		}),
		CapacityRemaining: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "CapacityRemaining",
			Help:      "CapacityRemaining",
		}),
		CapacityUsedNonDFS: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "CapacityUsedNonDFS",
			Help:      "CapacityUsedNonDFS",
		}),
		BlocksTotal: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "BlocksTotal",
			Help:      "BlocksTotal",
		}),
		FilesTotal: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "FilesTotal",
			Help:      "FilesTotal",
		}),
		CorruptBlocks: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "CorruptBlocks",
			Help:      "CorruptBlocks",
		}),
		ExcessBlocks: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "ExcessBlocks",
			Help:      "ExcessBlocks",
		}),
		StaleDataNodes: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "StaleDataNodes",
			Help:      "StaleDataNodes",
		}),
		pnGcCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "ParNew_CollectionCount",
			Help:      "ParNew GC Count",
		}),
		pnGcTime: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "ParNew_CollectionTime",
			Help:      "ParNew GC Time",
		}),
		cmsGcCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "ConcurrentMarkSweep_CollectionCount",
			Help:      "ConcurrentMarkSweep GC Count",
		}),
		cmsGcTime: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "ConcurrentMarkSweep_CollectionTime",
			Help:      "ConcurrentMarkSweep GC Time",
		}),
		heapMemoryUsageCommitted: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "heapMemoryUsageCommitted",
			Help:      "heapMemoryUsageCommitted",
		}),
		heapMemoryUsageInit: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "heapMemoryUsageInit",
			Help:      "heapMemoryUsageInit",
		}),
		heapMemoryUsageMax: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "heapMemoryUsageMax",
			Help:      "heapMemoryUsageMax",
		}),
		heapMemoryUsageUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "heapMemoryUsageUsed",
			Help:      "heapMemoryUsageUsed",
		}),
		lastHATransitionTime: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "lastHATransitionTime",
			Help:      "last HA Transition Time",
		}),
		state: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "state",
			Help:      "Current namenode state, 1 if active 0 if standby",
		}),
		isActive: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "isActive",
			Help:      "isActive",
		}),
	}
}

// Describe implements the prometheus.Collector interface.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.MissingBlocks.Describe(ch)
	e.UnderReplicatedBlocks.Describe(ch)
	e.CapacityTotal.Describe(ch)
	e.CapacityUsed.Describe(ch)
	e.CapacityRemaining.Describe(ch)
	e.CapacityUsedNonDFS.Describe(ch)
	e.BlocksTotal.Describe(ch)
	e.FilesTotal.Describe(ch)
	e.CorruptBlocks.Describe(ch)
	e.ExcessBlocks.Describe(ch)
	e.StaleDataNodes.Describe(ch)
	e.pnGcCount.Describe(ch)
	e.pnGcTime.Describe(ch)
	e.cmsGcCount.Describe(ch)
	e.cmsGcTime.Describe(ch)
	e.heapMemoryUsageCommitted.Describe(ch)
	e.heapMemoryUsageInit.Describe(ch)
	e.heapMemoryUsageMax.Describe(ch)
	e.heapMemoryUsageUsed.Describe(ch)
	e.lastHATransitionTime.Describe(ch)
	e.state.Describe(ch)
	e.isActive.Describe(ch)
}

// Collect implements the prometheus.Collector interface.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	resp, err := http.Get(e.url)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	var f interface{}
	err = json.Unmarshal(data, &f)
	if err != nil {
		log.Error(err)
	}
	// {"beans":[{"name":"Hadoop:service=NameNode,name=FSNamesystem", ...}, {"name":"java.lang:type=MemoryPool,name=Code Cache", ...}, ...]}
	m := f.(map[string]interface{})
	// [{"name":"Hadoop:service=NameNode,name=FSNamesystem", ...}, {"name":"java.lang:type=MemoryPool,name=Code Cache", ...}, ...]
	var nameList = m["beans"].([]interface{})
	for _, nameData := range nameList {
		nameDataMap := nameData.(map[string]interface{})
		/*
		   {
		       "name" : "Hadoop:service=NameNode,name=FSNamesystem",
		       "modelerType" : "FSNamesystem",
		       "tag.Context" : "dfs",
		       "tag.HAState" : "active",
		       "tag.TotalSyncTimes" : "23 6 ",
		       "tag.Hostname" : "CNHORTO7502.line.ism",
		       "MissingBlocks" : 0,
		       "MissingReplOneBlocks" : 0,
		       "ExpiredHeartbeats" : 0,
		       "TransactionsSinceLastCheckpoint" : 2007,
		       "TransactionsSinceLastLogRoll" : 7,
		       "LastWrittenTransactionId" : 172706,
		       "LastCheckpointTime" : 1456089173101,
		       "CapacityTotal" : 307099828224,
		       "CapacityTotalGB" : 286.0,
		       "CapacityUsed" : 1471291392,
		       "CapacityUsedGB" : 1.0,
		       "CapacityRemaining" : 279994568704,
		       "CapacityRemainingGB" : 261.0,
		       "CapacityUsedNonDFS" : 25633968128,
		       "TotalLoad" : 6,
		       "SnapshottableDirectories" : 0,
		       "Snapshots" : 0,
		       "LockQueueLength" : 0,
		       "BlocksTotal" : 67,
		       "NumFilesUnderConstruction" : 0,
		       "NumActiveClients" : 0,
		       "FilesTotal" : 184,
		       "PendingReplicationBlocks" : 0,
		       "UnderReplicatedBlocks" : 0,
		       "CorruptBlocks" : 0,
		       "ScheduledReplicationBlocks" : 0,
		       "PendingDeletionBlocks" : 0,
		       "ExcessBlocks" : 0,
		       "PostponedMisreplicatedBlocks" : 0,
		       "PendingDataNodeMessageCount" : 0,
		       "MillisSinceLastLoadedEdits" : 0,
		       "BlockCapacity" : 2097152,
		       "StaleDataNodes" : 0,
		       "TotalFiles" : 184,
		       "TotalSyncCount" : 7
		   }
		*/
		if nameDataMap["name"] == "Hadoop:service=NameNode,name=FSNamesystem" {
			e.MissingBlocks.Set(nameDataMap["MissingBlocks"].(float64))
			e.UnderReplicatedBlocks.Set(nameDataMap["UnderReplicatedBlocks"].(float64))
			e.CapacityTotal.Set(nameDataMap["CapacityTotal"].(float64))
			e.CapacityUsed.Set(nameDataMap["CapacityUsed"].(float64))
			e.CapacityRemaining.Set(nameDataMap["CapacityRemaining"].(float64))
			e.CapacityUsedNonDFS.Set(nameDataMap["CapacityUsedNonDFS"].(float64))
			e.BlocksTotal.Set(nameDataMap["BlocksTotal"].(float64))
			e.FilesTotal.Set(nameDataMap["FilesTotal"].(float64))
			e.CorruptBlocks.Set(nameDataMap["CorruptBlocks"].(float64))
			e.ExcessBlocks.Set(nameDataMap["ExcessBlocks"].(float64))
			e.StaleDataNodes.Set(nameDataMap["StaleDataNodes"].(float64))
		}
		/*
		   {
		       "name" : "Hadoop:service=NameNode,name=NameNodeStatus",
		       "modelerType" : "org.apache.hadoop.hdfs.server.namenode.NameNode",
		       "SecurityEnabled" : false,
		       "NNRole" : "NameNode",
		       "HostAndPort" : "namenode1.hdfs.tamr:50071",
		       "LastHATransitionTime" : 1484149009998,
		       "State" : "active"
		   }
		*/
		if nameDataMap["name"] == "Hadoop:service=NameNode,name=NameNodeStatus" {
			if nameDataMap["State"] == "active" {
				e.state.Set(1)
			} else {
				e.state.Set(0)
			}
			e.lastHATransitionTime.Set(nameDataMap["LastHATransitionTime"].(float64))
		}
		if nameDataMap["name"] == "java.lang:type=GarbageCollector,name=ParNew" {
			e.pnGcCount.Set(nameDataMap["CollectionCount"].(float64))
			e.pnGcTime.Set(nameDataMap["CollectionTime"].(float64))
		}
		if nameDataMap["name"] == "java.lang:type=GarbageCollector,name=ConcurrentMarkSweep" {
			e.cmsGcCount.Set(nameDataMap["CollectionCount"].(float64))
			e.cmsGcTime.Set(nameDataMap["CollectionTime"].(float64))
		}
		/*
		   "name" : "java.lang:type=Memory",
		   "modelerType" : "sun.management.MemoryImpl",
		   "HeapMemoryUsage" : {
		       "committed" : 1060372480,
		       "init" : 1073741824,
		       "max" : 1060372480,
		       "used" : 124571464
		   },
		*/
		if nameDataMap["name"] == "java.lang:type=Memory" {
			heapMemoryUsage := nameDataMap["HeapMemoryUsage"].(map[string]interface{})
			e.heapMemoryUsageCommitted.Set(heapMemoryUsage["committed"].(float64))
			e.heapMemoryUsageInit.Set(heapMemoryUsage["init"].(float64))
			e.heapMemoryUsageMax.Set(heapMemoryUsage["max"].(float64))
			e.heapMemoryUsageUsed.Set(heapMemoryUsage["used"].(float64))
		}

		if nameDataMap["name"] == "Hadoop:service=NameNode,name=FSNamesystem" {
			if nameDataMap["tag.HAState"] == "active" {
				e.isActive.Set(1)
			} else {
				e.isActive.Set(0)
			}
		}

	}
	e.MissingBlocks.Collect(ch)
	e.UnderReplicatedBlocks.Collect(ch)
	e.CapacityTotal.Collect(ch)
	e.CapacityUsed.Collect(ch)
	e.CapacityRemaining.Collect(ch)
	e.CapacityUsedNonDFS.Collect(ch)
	e.BlocksTotal.Collect(ch)
	e.FilesTotal.Collect(ch)
	e.CorruptBlocks.Collect(ch)
	e.ExcessBlocks.Collect(ch)
	e.StaleDataNodes.Collect(ch)
	e.pnGcCount.Collect(ch)
	e.pnGcTime.Collect(ch)
	e.cmsGcCount.Collect(ch)
	e.cmsGcTime.Collect(ch)
	e.heapMemoryUsageCommitted.Collect(ch)
	e.heapMemoryUsageInit.Collect(ch)
	e.heapMemoryUsageMax.Collect(ch)
	e.heapMemoryUsageUsed.Collect(ch)
	e.lastHATransitionTime.Collect(ch)
	e.state.Collect(ch)
	e.isActive.Collect(ch)
}

func main() {
	flag.Parse()

	exporter := NewExporter(*namenodeJmxUrl)
	prometheus.MustRegister(exporter)

	log.Printf("Starting Server: %s", *listenAddress)
	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
        <head><title>NameNode Exporter</title></head>
        <body>
        <h1>NameNode Exporter</h1>
        <p><a href="` + *metricsPath + `">Metrics</a></p>
        </body>
        </html>`))
	})
	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}
