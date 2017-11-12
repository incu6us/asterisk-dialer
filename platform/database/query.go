package database

var queries = make(map[int]string)

const (
    MsisdnAll                      = iota
    MsisdnInProgress
    MsisdnInProgressWithPagination
    MsisdnInProgressCount
)

const (
    InProgressSuffix = "status = '' OR status = 'progress' OR status = 'recall'"
)

func GetQuery(key int) string {
    queries[MsisdnInProgressCount] = `SELECT count(l.*) FROM msisdn_lists l, msisdn_priorities p WHERE l.id = p.msisdn_id AND (` + InProgressSuffix + `)`
    queries[MsisdnInProgress] = `SELECT l.*, p.* FROM msisdn_lists l, msisdn_priorities p WHERE l.id = p.msisdn_id AND (` + InProgressSuffix + `) ORDER BY %s %s, l.id`
    queries[MsisdnInProgressWithPagination] = `SELECT l.*, p.* FROM msisdn_lists l, msisdn_priorities p WHERE l.id = p.msisdn_id AND (` + InProgressSuffix + `) ORDER BY %s %s, l.id LIMIT %d OFFSET %d`
    return queries[key]
}
