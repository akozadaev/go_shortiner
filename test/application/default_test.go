package application

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"go.uber.org/fx/fxtest"
)

func TestAppStartsAndHandlesRoot(t *testing.T) {
	app := fxtest.New(t,
		Module,
	)

	app.RequireStart()

	time.Sleep(500 * time.Millisecond)

	resp, err := http.Get(fmt.Sprintf("http://localhost%s/", port))
	if err != nil {
		t.Fatalf("Ошибка при GET /: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус 200 OK, получено: %d", resp.StatusCode)
	}

	app.RequireStop()
}
