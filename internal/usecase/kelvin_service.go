package usecase

type kelvinService struct {
}

func NewKelvinService() *kelvinService {
	return &kelvinService{}
}

func (k *kelvinService) GetKelvin(tempC float64) float64 {
	return tempC + 273
}
