package node

// Node represents a machine in the cluster that can run tasks. It tracks the
// machine's resources, both total and allocated, to help the scheduler make
// decisions about task placement.
type Node struct {
	// Name is the identifier used to identify an individual node.
	Name string

	// Ip is the IP address of the node in the cluster.
	Ip string

	// Cores is the total number of CPU cores available on the node.
	Cores int

	// Memory is the total amount of memory available on the node.
	Memory int

	// MemoryAllocated is the amount of memory currently allocated
	// to tasks.
	MemoryAllocated int

	// Disk is the total amount of disk space available on the node.
	Disk int

	// DiskAllocated is the amount of disk space currently allocated
	// to tasks.
	DiskAllocated int

	// Role defines the role of the node in the cluster. Worker-Node,
	// Manager-Node
	Role string

	// TaskCount is the total number of tasks running on the node at
	// any given time.
	TaskCount int
}
