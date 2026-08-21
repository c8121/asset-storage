package metadata_db

type Finder interface {
	//Find searches in the metadata database
	Find(query any) (ScoredIdMap, error)
}
