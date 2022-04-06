package sources

import (
	"github.com/isacikgoz/morph/models"
)

type Source interface {
	Migrations() (migrations []*models.Migration)
}
