=====================================================================
 Block App
=====================================================================

:Source: https://github.com/frankofranko/ethApp

About
=====

Block App provides 2 functionalities:

1. APIs to query block and transaction info.

2. Block indexer

Setup
============

      .. code-block:: bash
      
         make setup


API server
============

Start server

.. code-block:: bash

     make start_server

Sample API requests

.. code-block:: bash

     curl http://localhost:8080/blocks?limit=5
     curl http://localhost:8080/blocks/18338591
     curl http://localhost:8080/transaction/0x1da684d542bd52aca3676217b536e9a4508dd0741e5ba27c17ec4714e1cef68a


Indexer
============

Indexer is a worker pool handling the unprocessed block number > $(block_num).

Start indexer:

Please replace the block_num to a reasonable number if you don't want to start with a big overhead.

.. code-block:: bash

   make start_indexer block_num=18338591 worker_num=10

Check DB

.. code-block:: bash
      
      make connect_dbshell

.. code-block:: MySQL

      select * from block limit 10;
