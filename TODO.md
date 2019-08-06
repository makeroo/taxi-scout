JS:
revert no submit decision: use preventDefault instead so that a user can submit by pressig enter

login / invitation:
* handle current state: automatically redirect to homepage if already signedin

* invitations: unit tests both on dao and rest components
  check "account w/o password workflow":
   - invitation => new account
   - cookie expires (the account has no password)
   - login cant be completed
   - select forgot password (implement it with a new invitation:YES
     because the account already exists and is member of a group!)

forgot password:
* handle send invitation response: show check your email msg or unknown email error

scouts editing:
 * obsolete initial state: discard scout property
 * selected group is missing: when adding a scout, scout group must be specified
   solution 1: support 1 scout group per account at most
   solution 2: evolve ui, add group selector (if needed)


async actions other than initial setup (that is "get"):
* choose an approach between "full redux" and "component logic"
 PRO full redux: testable
                 logic moved in async action
 PRO component logic: all code in one place


componentDidMount is called just once or everytime it is rendered?
 or: just the first time the router "activate" it or everytime it is "activated"?


when role changes to taxi, hide scouts in pickup summary?
  that is not true: what if two tutor of the same scouts
  pick some of them?

pickup summary: hide scouts that not partecipate

include google material icons as (npm?) dependency

study: websockets both react and go sides

test unit
