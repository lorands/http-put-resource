# http-put-resource

Concourse resource to be used to PUT a (filtered) content of folder to a HTTP endpoint with PUT.

This works pretty well with Nexus RAW repositories.


## Source Configuration

* `url`: *Required*. The http(s) url endpoint.
* `username`: *Optional* The username used to authenticate.
* `password`: *Optional*. The password used to authenticate.
* `verbose`: *Optional*. True to write intensive log.

## Check

No effect.

## Get

No effect.

## Put

Puts the provided folder to http. 

* `from`: *Required* The relative path to start watching for files
* `from-re-filter`: *Optional* Inclusive regesx filter to be used. Each path element is matched aginst it.
* `to`: *Required* The subpath on HTTP end. Can contain env variables. Look at the example.

## Pipeline example

```yaml
---
resource_types:
  - name: http-resource
    type: docker-image
    source:
      repository: lorands/http-put-resource
resources:
  - name: http-resource-nx-raw
    type: http-resource
    source:
      url: https://mynexus.example.com/repository/raw1
      username: myUser
      password: myPass
      verbose: true
jobs:
  - name: publish
    plan:
    - get: source-code
      trigger: true
    - get: version
      params:
        bump: patch
    - put: version
      params: {file: version/number}
    - task: publish
      file: source-code/ci/tasks/publish.yml
      ensure:
        put: http-resource-nx-raw
        params:
          from: build-result/
          from-re-filter: [^/]+/reports/.*
          to: builds/$BUILD_PIPELINE_NAME/$BUILD_JOB_NAME/$BUILD_NAME
```



