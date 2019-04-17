package perflog

func NewStore() *Store {
	return new(Store)
}

type Store struct {
	projects []*Project
}

func (s *Store) AddProject(id string, name string) *Project {
	p := &Project{
		Id:           id,
		Name:         name,
		Versions:     []Version{},
		Benchmarks:   []Benchmark{},
		Measurements: []Measurement{},
	}

	s.projects = append(s.projects, p)

	return p
}

func (s *Store) GetProjects() []*Project {
	return s.projects
}

type Project struct {
	Id           string
	Name         string
	Versions     []Version
	Benchmarks   []Benchmark
	Measurements []Measurement
}

func (a *Project) AddVersion(id string) {
	a.Versions = append(a.Versions, Version{Id: id, Name: id})
}

func (a *Project) AddBenchmark(id string) {
	a.Benchmarks = append(a.Benchmarks, Benchmark{Id: id, Name: id})
}

func (a *Project) AddMeasurement(versionId string, benchmarkId string, ops int) {
	a.Measurements = append(a.Measurements, Measurement{
		VersionId:           versionId,
		BenchmarkId:         benchmarkId,
		OperationsPerSecond: ops,
	})
}

type Version struct {
	Id   string
	Name string
}

type Benchmark struct {
	Id   string
	Name string
}

type Measurement struct {
	Id                  string
	Name                string
	VersionId           string
	BenchmarkId         string
	OperationsPerSecond int
}
