# ElasticSearch objects Operator

## Rationale
ElasticSearch operator is Kubernetes operator, which manages ES cluster indexes and templates.

## Indexes
Operator has folowing logic while managing ES indexes:
1. Create logic:

    1.1. If index already exists or there is an error in index settings, operator would report error in `status.latest_error` and would not perform any other operations upon that index. Only thing you could perform is to start from scratch: delete CRD, fix issues and try again.

    1.2. Upon index successful creation you would be able to manage that index via CRD.

2. Update logic:

    2.1. NOTE: If you change index name, operator would create new index according to last known configuration. Previously created index would not be deleted.

    2.2. If index update failed (due to incorrect values), error would be reported in `status.latest_error`. Please fix issue and update would be tried again.

3. Upon delete:
*NOTE*: CRD would be deleted from K8S in any case.

    3.1. If index have not been `aknowledged` from ES as successfully created or updated, operator would not try to delete it from ES.

    3.2. If CRD doesn't have `drop_on_delete` flag set, we would only delete CRD, but not ES index itself.

    3.3. If CRD has `drop_on_delete` flag set, we would also try to delete index. If error occures, it would be reported to operator logs.

## Templates
Operator has folowing logic while managing ES templates:
1. Create logic:

    1.1. If template alredy exists or there is an error in template settings, operator would report error in `status.latest_error` and would not perform any other operations upon that template. Only thing you could perform is to start from scratch: delete CRD, fix issues and try again.

    1.2. Upon template successful creation you would be able to manage that template via CRD.

2. Update logic:

    2.1. NOTE: If you change template name, operator would create new template according to last known configuration. Previously created template would not be deleted.

    2.2. If template update failed (due to incorrect values), error would be reported in `status.latest_error`. Please fix issue and update would be tried again.

3. Upon delete:
*NOTE*: CRD would be deleted from K8S in any case.

    3.1. If template have not been `aknowledged` from ES as successfully created or updated, operator would not try to delete it from ES.

    3.2. If CRD doesn't have `drop_on_delete` flag set, we would only delete CRD, but not ES template itself.

    3.3. If CRD has `drop_on_delete` flag set, we would also try to delete template. If error occures, it would be reported to operator logs.

Documentation is available [here](https://elasticsearch-objects-operator.readthedocs.io/en/latest/).
