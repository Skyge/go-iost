class Account {
    constructor() {

    }
    init() {

    }
    _saveAccount(account) {
        storage.put("auth", JSON.stringify(account), account.id)
    }
    _loadAccount(id) {
        let a = storage.get("auth", id);
        return JSON.parse(a)
    }
    static _find(items, name) {
        for (let i = 0; i < items.length(); i ++) {
            if (items[i].id === name ) {
                return i
            }
        }
        return -1
    }

    _hasAccount(id) {
        return storage.has("auth", id)
    }

    _ra(id) {
        if (!BlockChain.requireAuth(id, "owner")) {
            throw "require auth failed"
        }
    }

    SignUp(id, owner, active) {
        if (this._hasAccount(id)) {
            throw "id existed > " + id
        }
        let account = {};
        account.id = id;
        account.permissions = {};
        account.permissions.active = {
            name: "active",
            groups: [],
            items: [{
                id: active,
                is_key_pair: true,
                weight: 1,
            }],
            threshold: 1,
        };
        account.permissions.owner = {
            name: "owner",
            groups: [],
            items: [{
                id: owner,
                is_key_pair: true,
                weight: 1,
            }],
            threshold: 1,
        };
        this._saveAccount(account)
    }
    AddPermission(id, perm, thres) {
        this._ra(id);
        let acc = this._loadAccount(id);
        if (acc.permissions[perm] !== undefined) {
            throw "permission already exist"
        }
        acc.permissions[perm] = {
            name: perm,
            groups: [],
            items: [],
            threshold: thres,
        };
        this._saveAccount(acc)
    }
    DropPermission(id, perm) {
        this._ra(id);
        let acc = this._loadAccount(id);
        acc.permissions[perm] = undefined;
        this._saveAccount(acc)
    }
    AssignPermission(id, perm, un, weight) {
        this._ra(id);
        let acc = this._loadAccount(id);
        const index = Account._find(acc.permissions[perm].items, un);
        if (index < 0) {
            const len = un.indexOf("@");
            if (len < 0 && un.startsWith("IOST") ){
                acc.permissions[perm].items.push({
                    id : un,
                    is_key_pair: true,
                    weight: weight
                });
            } else {
                acc.permissions[perm].items.push({
                    id : un.substring(0, len),
                    permission: un.substring(len, un.length()),
                    is_key_pair: false,
                    weight: weight
                });
            }
        } else {
            acc.permissions[perm].items[index].weight = weight
        }
        this._saveAccount(acc)
    }
    RevokePermission(id, perm, un) {
        this._ra(id);
        let acc = this._loadAccount(id);
        const index = Account._find(acc.permissions[perm].items, un);
        if (index < 0) {
            throw "item not found"
        } else {
            acc.permissions[perm].items.splice(index, 1)
        }
        this._saveAccount(acc)
    }
    AddGroup(id, grp) {
        this._ra(id);
        let acc = this._loadAccount(id);
        if (acc.groups[grp] !== undefined) {
            throw "group already exist"
        }
        acc.groups[grp] = {
            name : grp,
            items : [],
        };
        this._saveAccount(acc)
    }
    DropGroup(id, group) {
        this._ra(id);
        let acc = this._loadAccount(id);
        acc.permissions[group] = undefined;
        this._saveAccount(acc)
    }
    AssignGroup(id, group, un, weight) {
        this._ra(id);
        let acc = this._loadAccount(id);
        const index = Account._find(acc.permissions[group].items, un);
        if (index < 0) {
            let len = un.indexOf("@");
            if (len < 0 && un.startsWith("IOST")) {
                acc.groups[group].items.push({
                    id: un,
                    is_key_pair: true,
                    weight: weight
                });
            } else {
                acc.groups[group].items.push({
                    id: un.substring(0, len),
                    permission: un.substring(len, un.length()),
                    is_key_pair: false,
                    weight: weight
                });
            }
        } else {
            acc.groups[group].items[index].weight = weight
        }

        this._saveAccount(acc)
    }
    RevokeGroup(id, grp, un) {
        this._ra(id);
        let acc = this._loadAccount(id);
        const index = Account._find(acc.permissions[grp].items, un);
        if (index < 0) {
            throw "item not found"
        } else {
            acc.permissions[grp].items.splice(index, 1)
        }
        this._saveAccount(acc)
    }
    AssignPermissionToGroup(id, perm, group) {
        this._ra(id);
        let acc = this._loadAccount(id);
        acc.permissions[perm].groups.push(group);
        this._saveAccount(acc)
    }
    RevokePermissionInGroup(id, perm, group) {
        this._ra(id);
        let acc = this._loadAccount(id);
        let index = acc.permissions[perm].groups.indexOf(group);
        if (index > -1) {
            acc.permissions[perm].groups.splice(index, 1);
        }
        this._saveAccount(acc)
    }
}

module.exports = Account;