package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Funkcja New nie może zwracać nil!")
	} else {
		tracer.Trace("Witamy w pakiecie trace")
		if buf.String() != "Witamy w pakiecie trace\n" {
			t.Errorf("Metoda Trace nie powinna wyświetlać '%s", buf.String())
		}
	}
}
func TestOff(t *testing.T) {
	silentTracer := Off()
	silentTracer.Trace("cośtam")
}
