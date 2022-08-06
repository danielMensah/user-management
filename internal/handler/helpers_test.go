package handler

import (
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func teardown(mt *mtest.T) {
	mt.ClearMockResponses()
	mt.ClearCollections()
	mt.ClearFailPoints()
	mt.ClearEvents()
}
