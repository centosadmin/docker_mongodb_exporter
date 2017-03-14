package collector_mongos

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	tcMallocCurrentBytes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "tcmalloc",
		Name:      "bytes",
		Help:      "The current number of bytes used by TCMalloc",
	}, []string{"type"})
	tcMallocCurrentPageHeapBytes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "tcmalloc",
		Name:      "page_heap_bytes",
		Help:      "The current number of page heap bytes used by TCMalloc",
	}, []string{"type"})
	tcMallocThreadCacheBytes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "tcmalloc",
		Name:      "thread_cache_bytes",
		Help:      "The number of thread cache bytes used by TCMalloc",
	}, []string{"type"})
	tcMallocAggresiveMemDecommit = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "tcmalloc",
		Name:      "aggressive_memory_decommit",
		Help:      "The current number of aggressive memory decommits by TCMalloc",
	})
)

type TCMallocGenericStats struct {
	CurrrentAllocatedBytes float64 `bson:"current_allocated_bytes,omitempty"`
	HeapSizeBytes          float64 `bson:"heap_size,omitempty"`
}

func (t *TCMallocGenericStats) Export(ch chan<- prometheus.Metric) {
	tcMallocCurrentBytes.WithLabelValues("allocated").Set(t.CurrrentAllocatedBytes)
	tcMallocCurrentBytes.WithLabelValues("heap").Set(t.HeapSizeBytes)
}

type TCMallocTCMallocStats struct {
	PageHeapFreeBytes            float64 `bson:"pageheap_free_bytes,omitempty"`
	PageHeapUnmappedBytes        float64 `bson:"pageheap_unmapped_bytes,omitempty"`
	MaxTotalThreadCacheBytes     float64 `bson:"max_total_thread_cache_bytes,omitempty"`
	CurrentTotalThreadCacheBytes float64 `bson:"current_total_thread_cache_bytes,omitempty"`
	TotalFreeBytes               float64 `bson:"total_free_bytes,omitempty"`
	CentralCacheFreeBytes        float64 `bson:"central_cache_free_bytes,omitempty"`
	TransferCacheFreeBytes       float64 `bson:"transfer_cache_free_bytes,omitempty"`
	ThreadCacheFreeBytes         float64 `bson:"thread_cache_free_bytes,omitempty"`
	AggressiveMemoryDecommit     float64 `bson:"aggressive_memory_decommit,omitempty"`
}

func (t *TCMallocTCMallocStats) Export(ch chan<- prometheus.Metric) {
	tcMallocCurrentBytes.WithLabelValues("free").Set(t.TotalFreeBytes)
	tcMallocCurrentBytes.WithLabelValues("cache_free").Set(t.CentralCacheFreeBytes)
	tcMallocCurrentBytes.WithLabelValues("transfer_cache_free").Set(t.TransferCacheFreeBytes)
	tcMallocCurrentPageHeapBytes.WithLabelValues("free").Set(t.PageHeapFreeBytes)
	tcMallocCurrentPageHeapBytes.WithLabelValues("unmapped").Set(t.PageHeapUnmappedBytes)
	tcMallocThreadCacheBytes.WithLabelValues("max").Set(t.MaxTotalThreadCacheBytes)
	tcMallocThreadCacheBytes.WithLabelValues("total").Set(t.CurrentTotalThreadCacheBytes)
}

type TCMallocStats struct {
	Generic  *TCMallocGenericStats  `bson:"generic,omitempty"`
	TCMalloc *TCMallocTCMallocStats `bson:"tcmalloc,omitempty"`
}

func (stats *TCMallocStats) Describe(ch chan<- *prometheus.Desc) {
	tcMallocCurrentBytes.Describe(ch)
	tcMallocCurrentPageHeapBytes.Describe(ch)
	tcMallocThreadCacheBytes.Describe(ch)
	tcMallocAggresiveMemDecommit.Describe(ch)
}

func (stats *TCMallocStats) Export(ch chan<- prometheus.Metric) {
	if stats.Generic != nil {
		stats.Generic.Export(ch)
	}
	if stats.TCMalloc != nil {
		stats.TCMalloc.Export(ch)
		tcMallocCurrentPageHeapBytes.Collect(ch)
		tcMallocAggresiveMemDecommit.Collect(ch)
		tcMallocThreadCacheBytes.Collect(ch)
	}
	tcMallocCurrentBytes.Collect(ch)
}
