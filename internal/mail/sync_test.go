package mail_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tauraamui/maildew/internal/mail"
)

func TestResolveNewAndDeletedUIDsFromSetsSuccess(t *testing.T) {
	is := is.New(t)

	local := []uint32{122, 123, 124, 125, 128}
	remote := []uint32{122, 124, 125, 128, 129}

	new, missing := mail.ResolveAddedAndRemoved(local, remote)

	is.Equal(new, []uint32{129})
	is.Equal(missing, []uint32{123})
}
