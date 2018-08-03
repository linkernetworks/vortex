#!/usr/bin/env bats

load init

@test "Create Storage" {
    http -v --check-status 127.0.0.1:7890/v1/storage < storage.json
    [ $? = 0 ]
}

@test "List Storage" {
    run bash -c 'http http://127.0.0.1:7890/v1/storage/ 2>/dev/null | jq -r ".[0].id"'
    id=${output}
    run kubectl -n vortex get sc nfs-storageclass-${id} -o jsonpath="{.provisioner}"
    [ "$output" = "nfs-provisioner-${id}" ]

    NEXT_WAIT_TIME=0
    WAIT_LIMIT=40
    deploymentName="nfs-provisioner-${id}"
    until kubectl -n vortex get deployment ${deploymentName} -o jsonpath="{.status.readyReplicas}" | grep "1" || [ $NEXT_WAIT_TIME -eq $WAIT_LIMIT ]; do
       sleep 2
       kubectl -n vortex get deployment ${deploymentName}
       NEXT_WAIT_TIME=$((NEXT_WAIT_TIME+ 1))
    done
    [ $NEXT_WAIT_TIME != $WAIT_LIMIT ]
}

@test "Create Volume" {
    http -v --check-status 127.0.0.1:7890/v1/volume < volume.json
    [ $? = 0 ]
}

@test "List Volume" {
    run bash -c 'http http://127.0.0.1:7890/v1/volume/ 2>/dev/null | jq -r ".[0].id"'
    id=${output}
    NEXT_WAIT_TIME=0
    WAIT_LIMIT=40
    pvcName="pvc-${id}"
    until kubectl get pvc ${pvcName} -o jsonpath="{.status.phase}" | grep "Bound" || [ $NEXT_WAIT_TIME -eq $WAIT_LIMIT ]; do
       sleep 2
       kubectl get pvc ${pvcName}
       NEXT_WAIT_TIME=$((NEXT_WAIT_TIME+ 1))
    done
    [ $NEXT_WAIT_TIME != $WAIT_LIMIT ]
}

@test "Create Pod" {
    http -v --check-status 127.0.0.1:7890/v1/pods < pod.json
    [ $? = 0 ]
}

@test "List Pod" {
   run bash -c "http http://127.0.0.1:7890/v1/pods/ 2>/dev/null | jq -r '.[] | select(.name == \"${podName}\").name'"
   [ "$output" = "${podName}" ]
   [ $status = 0 ]

   NEXT_WAIT_TIME=0
   WAIT_LIMIT=40
   until kubectl get pods ${podName} -o jsonpath="{.status.phase}" | grep "Running" || [ $NEXT_WAIT_TIME -eq $WAIT_LIMIT ]; do
      sleep 2
      kubectl get pods ${podName}
      NEXT_WAIT_TIME=$((NEXT_WAIT_TIME+ 1))
   done
   [ $NEXT_WAIT_TIME != $WAIT_LIMIT ]
}

@test "Test NFS" {
    kubectl exec ${podName} touch /tmp/testing
    find /tmp/nfs | grep testing
    [ $? = 0 ]
}

@test "Delete Pod" {
    run bash -c 'http http://127.0.0.1:7890/v1/pods/ 2>/dev/null | jq -r ".[0].id"'
    run http DELETE http://127.0.0.1:7890/v1/pods/${output} 2>/dev/null
    [ $status = 0 ]
    NEXT_WAIT_TIME=0
    WAIT_LIMIT=100
    until kubectl get pods ${podName} 2>&1 | grep "No resources" || [ $NEXT_WAIT_TIME -eq $WAIT_LIMIT ]; do
      sleep 2
      kubectl get pods
      NEXT_WAIT_TIME=$((NEXT_WAIT_TIME+ 1))
   done
   [ $NEXT_WAIT_TIME != $WAIT_LIMIT ]
}

@test "Delete Volume" {
    run bash -c 'http http://127.0.0.1:7890/v1/volume/ 2>/dev/null | jq -r ".[0].id"'
    http -v --check-status DELETE http://127.0.0.1:7890/v1/volume/${output}
    [ $? = 0 ]
}

@test "Delete Storage" {
    run bash -c 'http http://127.0.0.1:7890/v1/storage/ 2>/dev/null | jq -r ".[0].id"'
    http -v --check-status DELETE http://127.0.0.1:7890/v1/storage/${output}
    [ $? = 0 ]
}
