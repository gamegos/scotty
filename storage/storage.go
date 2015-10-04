package storage

type driverFactory func(config map[string]interface{}) Storage

var drivers = make(map[string]driverFactory)

func Register(driverType string, factory driverFactory) {
	drivers[driverType] = factory
}

func Init(driverType string, options map[string]interface{}) Storage {
	factory, ok := drivers[driverType]
	if !ok {
		panic("storage:driver:" + driverType + " is unknown")
	}
	return factory(options)
}
