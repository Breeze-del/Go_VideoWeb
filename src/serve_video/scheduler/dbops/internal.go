package dbops

import "log"

func ReadVideoDeletionRecord(count int) ([]string, error) {
	stmOut, err := dbConn.Prepare("SELECT video_id FROM video_del_rec LIMIT ?")
	var ids []string
	if err != nil {
		return ids, err
	}
	rows, err := stmOut.Query(count)
	if err != nil {
		log.Printf("Query videoDeletion error %s", err)
		return ids, err
	}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return ids, err
		}
		ids = append(ids, id)
	}
	stmOut.Close()
	return ids, nil
}

func DelVideoDeletionRecoed(vid string) error {
	stmDel, err := dbConn.Prepare("DELETE FROM video_del_rec WHERE video_id =?")
	if err != nil {
		return err
	}
	_, err = stmDel.Exec(vid)
	if err != nil {
		return err
	}
	stmDel.Close()
	return nil
}
