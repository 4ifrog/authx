db.createUser({
    user: 'nobody',
    pwd:  'secrets',
    roles: [
        {
            role: 'dbAdmin',
            db:   'authx',
        },
        {
            role: 'readWrite',
            db:   'authx',
        },
    ]
});
