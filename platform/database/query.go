package database

var queries = make(map[int]string)

const (
    MsisdnAll                      = iota
    MsisdnInProgress
    MsisdnInProgressWithPagination
)

const (
    InProgressSuffix = "status = '' OR status = 'progress' OR status = 'recall'"
)

func GetQuery(key int) string {
    queries[MsisdnInProgress] = `SELECT l.*, p.* FROM msisdn_lists l, msisdn_priorities p WHERE l.id = p.msisdn_id AND (` + InProgressSuffix + `) ORDER BY %s %s`
    queries[MsisdnInProgressWithPagination] = `SELECT l.*, p.* FROM msisdn_lists l, msisdn_priorities p WHERE l.id = p.msisdn_id AND (` + InProgressSuffix + `) ORDER BY %s %s LIMIT %d OFFSET %d`
    return queries[key]
}
