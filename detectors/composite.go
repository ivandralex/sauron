package detectors

import "github.com/sauron/session"

//CompositeDetector uses all simple detectors to get label
type CompositeDetector struct {
	composition []Detector
}

//AddDetector add another detector to the composition
func (d *CompositeDetector) AddDetector(detector Detector) {
	d.composition = append(d.composition, detector)
}

//Init initializes composite detector
func (d *CompositeDetector) Init(configPath string) {}

//GetLabel returns the first non-zero label returned by composition of detectors
func (d *CompositeDetector) GetLabel(s *session.SessionData) int {
	var label int
	for _, detector := range d.composition {
		label = detector.GetLabel(s)

		if label != UnknownLabel {
			return label
		}
	}

	return UnknownLabel
}
