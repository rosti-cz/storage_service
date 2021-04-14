# Storage service

Storage service listens for message coming to NATS servers/clusters and created/updated/deletes databases based on them.


## Events

This service listens to following events:

    subject: admin.storages.{storage_type}.{server}.events
    {
        event_type: "deleted"
        db_id:      int
        db_name:    string
        username:   string
    }

    subject: admin.storages.{storage_type}.{server}.events
    {
        event_type: "created"
        db_name:    string
        db_id:      int.id
        username:   string.username
        password:   string
        extensions: string
    }

    subject: admin.storages.{storage_type}.{server}.events
    {
        event_type: "password_changed"
        db_name:    string
        db_id:      int
        username:   string
        password:   string
    }

This service also emits state messages

    subject: admin.storages.{storage_type}.{server}.states
    {
        db_id:   int
	    db_name: string
	    error:   bool
	    message: string
    }
    
