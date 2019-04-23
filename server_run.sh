url0=tcp://127.0.0.1:40890
url1=tcp://127.0.0.1:40891
url2=tcp://127.0.0.1:40892
url3=tcp://127.0.0.1:40893
./buspubsub server0 $url0 $url1 $url2 & server0=$!
./buspubsub server1 $url1 $url2 $url3 & server1=$!
./buspubsub server2 $url2 $url3 & server2=$!
./buspubsub server3 $url3 $url0 & server3=$!
sleep 5
kill $server0 $server1 $server2 $server3