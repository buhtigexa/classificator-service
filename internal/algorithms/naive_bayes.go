package algorithms

type Class struct {
	frequency float64
	terms     map[string]float64
}

type NaiveBayes struct {
	classes map[string]*Class
}

func NewNaiveBayes() *NaiveBayes {
	return &NaiveBayes{make(map[string]*Class)}
}

func (n *NaiveBayes) Train(w string, class string) {
	if _, exists := n.classes[class]; !exists {
		n.classes[class] = &Class{0.0, make(map[string]float64)}
	}
	n.classes[class].terms[w]++
	n.classes[class].frequency++
}
