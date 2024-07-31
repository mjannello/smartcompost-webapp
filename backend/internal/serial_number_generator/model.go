package serial_number_generator

type SerialNumberGenerator interface {
	Generate() (string, error)
}
