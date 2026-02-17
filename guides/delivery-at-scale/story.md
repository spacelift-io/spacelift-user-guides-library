# The Saturnhead Startup Chronicles: Delivery at Scale

_Orbit Labs isn't one stack anymore. It's a system — networking feeds into compute, compute feeds into applications, and everything needs to deploy in the right order. Saturnhead can't coordinate this manually. He needs orchestration._

---

## Mission 1: Launchpad (Speedrun)

_"Before you can orchestrate, you need something to orchestrate."_

Saturnhead needs the basics in place — fast. VCS connected, AWS wired up, a stack that does something real.

### Tasks

1. Connect a VCS provider
2. Create a repository with Terraform that provisions an S3 bucket (we provide the code)
3. Create an AWS integration with cross-account role assumption
4. Create a stack, attach the integration
5. Trigger a run, confirm it completes

### API Checks

- VCS integration exists
- AWS integration exists with role assumption
- Stack exists with AWS integration attached
- At least one run in `FINISHED` state with resources created

---

## Mission 2: One Leads, One Follows

_"The database has to exist before the app can connect to it."_

Saturnhead creates two stacks. The app stack depends on the database stack. When the database deploys, the app follows automatically.

### Tasks

1. Create a second stack: `database` (we provide: Terraform that creates an S3 bucket simulating a "data store")
2. Rename your original stack to `app`
3. Create a dependency: `app` depends on `database`
4. Trigger the `database` stack
5. Observe the `app` stack triggers automatically after `database` finishes

### API Checks

- Two stacks exist: `database` and `app`
- Dependency configured: `app` → `database`
- `database` run completes
- `app` run triggered automatically after `database` completes

### Learning Beat

Dependencies are explicit. No guessing, no race conditions.

---

## Mission 3: Pass the Data

_"The app doesn't just need the database to exist — it needs to know where it is."_

Stack dependencies can pass outputs from one stack as inputs to another. The database stack outputs a bucket name. The app stack consumes it.

### Tasks

1. Modify `database` Terraform to output the bucket name (we provide the code)
2. Modify `app` Terraform to accept a `data_bucket` input variable and use it (we provide the code)
3. Configure the dependency to wire `database` output → `app` input
4. Trigger `database`
5. Observe `app` receives the bucket name and uses it

### API Checks

- `database` stack has output configured
- Dependency has output-to-input mapping
- `app` run receives input value
- `app` run outputs show the consumed value

### Learning Beat

No more hardcoded values. No more "go look at the other stack and copy-paste." Data flows through the graph.

---

## Mission 4: Wait, Skip, Proceed

_"Not every push changes everything."_

When a dependency runs but produces no changes, the downstream stack doesn't need to run. Spacelift skips it. When changes are in progress, downstream stacks wait.

### Tasks

1. Trigger `database` with no code changes — observe a "no changes" run
2. Observe `app` is **skipped** (not triggered, not failed — just skipped)
3. Make a real change to `database` (e.g., add a tag)
4. Observe `app` enters **pending** state while `database` runs
5. After `database` completes, observe `app` proceeds

### API Checks

- `database` run with no changes detected
- `app` run skipped (not triggered)
- `database` run with changes
- `app` run in pending state during `database` execution
- `app` run proceeds after `database` completes

### Learning Beat

Smart orchestration. Don't run what doesn't need to run. Don't run out of order.

---

## Mission 5: The Promotion Gate

_"Staging has to pass before production can even try."_

Saturnhead wants a simple promotion pattern: push code, staging deploys automatically, production waits for human confirmation — but only if staging succeeded.

### Tasks

1. Create a new repo (or folder) with shared Terraform code (we provide: simple S3 bucket with `env` variable)
2. Create two stacks pointing at the same code:
   - `staging` with `TF_VAR_env=staging` and **autodeploy on**
   - `production` with `TF_VAR_env=production` and **autodeploy off**
3. Create a dependency: `production` depends on `staging`
4. Push a change — observe `staging` runs automatically
5. After `staging` succeeds, observe `production` is **unconfirmed**, waiting for manual approval
6. Confirm the `production` run — observe it applies
7. Now break the Terraform (we provide: syntax error), push again
8. Observe `staging` fails — and `production` **never runs**

### API Checks

- Two stacks: `staging` (autodeploy on), `production` (autodeploy off)
- Dependency: `production` → `staging`
- Push triggers `staging` automatically
- `staging` success → `production` unconfirmed, awaiting approval
- `production` confirmed → applies successfully
- `staging` failure → `production` not triggered

### Learning Beat

The dependency isn't just about order — it's about safety. A failed upstream blocks downstream. You don't ship broken code to production because staging is your gate.
