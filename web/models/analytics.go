package models

type (
	AnalyticsData struct {
		PlayCount  uint   `json:"playcount"`
		UserCount  uint   `json:"usercount"`
		MCVersion  string `json:"mcversion"`
		ModVersion string `json:"modversion"`
	}
)

func (m *Model) GetPlayCount() (c uint, err error) {
	err = m.db.Get(&c, "SELECT SUM(updated_count) FROM user")
	return
}

func (m *Model) GetUserCount() (c uint, err error) {
	err = m.db.Get(&c, "SELECT count(uuid) FROM user")
	return
}

func (m *Model) GetMostPlayedMCVersion() (v string, err error) {
	err = m.db.Get(&v, "SELECT version_mc from user__version_mc__version_mod group by version_mc having count(*) = (select max(cnt) from (select count(*) as cnt from user__version_mc__version_mod group by version_mc) mc)")
	return
}

func (m *Model) GetMostPlayedModVersion() (v string, err error) {
	err = m.db.Get(&v, "SELECT version_mod from user__version_mc__version_mod group by version_mod having count(*) = (select max(cnt) from (select count(*) as cnt from user__version_mc__version_mod group by version_mod) mc)")
	return
}

func (m *Model) GetAnalyticsData() (d AnalyticsData, err error) {
	a := new(AnalyticsData)
	if a.PlayCount, err = m.GetPlayCount(); err != nil {
		return
	}
	if a.UserCount, err = m.GetUserCount(); err != nil {
		return
	}
	if a.MCVersion, err = m.GetMostPlayedMCVersion(); err != nil {
		return
	}
	if a.ModVersion, err = m.GetMostPlayedModVersion(); err != nil {
		return
	}
	return *a, nil
}
