# Gossip Package
## Contents
- `types.go`: type declarations for gossiper context and some basic gossiper methods definitions
- `gossiper.go`: implementation of gossiper functions
- `gossiper_object.go`: implementation of gossip object data type and signature verification methods for different type of gossip object, code in this file will be used by other entites in the system.

## types.go
- `Gossiper_context`: Gossiper context is an object that contains all the configuration and storage information about the gossiper
- `methods`: internal methods defined in this file includes save storage, load storage, wipe storage, store object, get object, isduplicate, hasPoM
## Gossiper.go
- `handlegossip`: gossip object handler
- `GossipData`: send gossip object to its connected gossipers
- `SendtoOwner`: send gossip object to its Owner (e.g. monitor)
- `ProcessValidObject`: valid Gossiper object handler, only invoked if the gossip object received has passed the signature verification process
- `ProcessDuplicateObject`: duplicate object handler, only invoked if the gossip object received has the same Gossip ID with a gossip object in the storage/cache, It will generate CONFLICT POM if the signature(s)/payload(s) are different
- `Process_TSS_Object`: Handlers for different types of gossip object with Threshold Signature Scheme 
- `PeriodicTasks` : only keeps the data for the current period (keep all Conflict PoMs), wipe every 3 Period 


