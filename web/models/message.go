package models

type (
	Message struct {
		UUID            string `json:"uuid" query:"id" validate:"required,mcuuid"`
		UserName        string `json:"username" query:"name" validate:"required"`
		IP              string `json:"ip"`
		VersionMod      string `json:"version_mod" query:"vmod" validate:"required" json:"-"`
		VersionModMC    string `json:"version_mod_mc" query:"vmodmc" validate:"required" json:"-"`
		VersionModForge string `json:"version_mod_forge" query:"vmodforge" validate:"required" json:"-"`
		VersionMC       string `json:"version_mc" query:"vmc" validate:"required" json:"-"`
		VersionForge    string `json:"version_forge" query:"vforge" validate:"required" json:"-"`
		Message         string `json:"message"`
		CreatedAt       string `json:"created_at"`
		UpdatedAt       string `json:"updated_at"`
		UpdatedCount    uint   `json:"updated_count"`
	}
)


func (m *Model) PutMessage(msg *Message) error {
	tx := m.db.MustBegin()
	tx.MustExec(`INSERT INTO user (uuid, username, ip, message) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE username = VALUES(username), ip = VALUES(ip), message = VALUES(message), updated_at = NOW(), updated_count = updated_count+1`, 
		msg.UUID,
		msg.UserName,
		msg.IP,
		msg.Message)
	tx.MustExec(`INSERT INTO user__version_mc__version_mod (uuid, version_mod, version_mod_mc, version_mod_forge, version_mc, version_forge) VALUES (?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE version_mod = VALUES(version_mod), version_mod_mc = VALUES(version_mod_mc), version_mod_forge = VALUES(version_mod_forge), version_mc = VALUES(version_mc), version_forge = (version_forge), updated_at = NOW(), updated_count = updated_count + 1`,
		msg.UUID,
		msg.VersionMod,
		msg.VersionModMC,
		msg.VersionModForge,
		msg.VersionMC,
		msg.VersionForge)
	return tx.Commit()
}
