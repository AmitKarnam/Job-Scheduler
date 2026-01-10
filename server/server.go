package server

const (
	defaultPort = 9001
)

// Server is the entrypoint of the job scheduler
type Server struct {
	port int
}

// Constructor like object used to return a server object
func NewServer(port int) *Server {
	if port == 0 {
		port = defaultPort
	}

	return &Server{
		port: port,
	}
}

// starting point of the server
func (s *Server) Start() error {
	// Start Database
	// 1. Check Connectivity to database.
	// 2. Check if the required tables are present
	// 3. If the AckLevel and ReadLevel are not present, then initialise the table and the value to 0
	// 4. if the Jobs tables is not present, initialise the table
	// Start Scheduler
	// 1. Fetch The AckLevel and ReadLevel.
	// 2. Create a in memory min-heap timer based on the AckLevel and ReadLevel
	// 3. Start a timer that takes care of AckLevel and Readlevel ( process that will populate the min heap )
	// 4. Start a timer that executes the min heap ( process that will execute the min heap jobs one by one )
	// Check if the system is NTP sync
	return nil
}

// stopping the server
func (s *Server) Stop() error {
	return nil
}
