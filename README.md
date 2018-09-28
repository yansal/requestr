# Requestr

## Design

* where should be the loop parameter? in the message? in the worker pool? in the receiver? (it's in the message for now)

## TODO

* Design remove
* Design pool healthcheck
* Use an RPC system for pool: interface could be something like

    type Pool interface {
        // These two are async and use internal broker
        Add(context.Context, queue string) error
        Remove(context.Context, queue string) error

        // Use an internal store?
        Healthcheck(context.Context) (..., error)
    }