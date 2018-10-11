package main

const user = `REPLACE INTO user (uuid, username, ip, version_mod, version_mod_mc, version_mod_forge, version_mc, version_forge) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

func (u *User) AddUser() {
    db.MustExec(user, u.UUID,
                      u.UserName,
                      u.IP, u.VersionMod,
                      u.VersionModMC,
                      u.VersionModForge,
                      u.VersionMC,
                      u.VersionForge)
}

