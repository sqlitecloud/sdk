-- /* launch n concurrent tasks, each task perform a set of read and write statements */

--set dbname=tester_nwriters3.sqlite
--set ntasks=20
--set ntimes=400
--set nread=2
--set nwrite=1
--set sleepms=1
--set beginstmt="BEGIN TRANSACTION"

CREATE DATABASE {{.dbname}} IF NOT EXISTS; USE DATABASE {{.dbname}};
--match type is ok
CREATE TABLE IF NOT EXISTS t1(idx INTEGER PRIMARY KEY, task INTEGER, loop INTEGER, write INTEGER); DELETE FROM t1;
--match type is ok

--loop t=1; t<=ntasks; t=t+1;
    --task t NAME readwrite
        --set initialsleep=t*sleepms
        --sleep initialsleep

        USE DATABASE {{.dbname}}

        -- /* each worker perform ntimes a set of nread reads and nwrite writes */
        --loop j=1; j<=ntimes; j=j+1;
            --loop r=1; r<=nread; r=r+1
                SELECT COUNT(*) FROM t1;
            --end
            --loop w=1; w<=nwrite; w=w+1;
                {{.beginstmt}}; INSERT INTO t1 (task, loop, write) VALUES ({{.t}}, {{.j}}, {{.w}}); COMMIT;
                -- // INSERT INTO t1 (task, loop, write) VALUES ({{.t}}, {{.j}}, {{.w}});
            --end

            --sleep sleepms
        --end

        UNUSE DATABASE
    --end
--end

--wait all

--loop i=1; i<=ntasks; i=i+1;
    --task i NAME exit
        --exit
    --end
--end

--wait all

SELECT COUNT(*) FROM t1;
--dump

UNUSE DATABASE;
--match type is ok

--sleep 100
LIST DATABASE CONNECTIONS {{.dbname}};

-- REMOVE DATABASE {{.dbname}};
-- --match type is ok
