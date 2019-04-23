url=tcp://127.0.0.1:40899
./buspubsub client $url Hello_World client & client=$! && sleep 1
./buspubsub server $url server0 & server0=$!
./buspubsub server $url server1 & server1=$!
./buspubsub server $url server2 & server2=$!
./buspubsub server $url server3 & server3=$!
sleep 1
kill $client $server0 $server1 $server2 $server3