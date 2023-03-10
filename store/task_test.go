package store

import (
	"context"
	"testing"

	"github.com/Satoshi-Tb/go_todo_app/clock"
	"github.com/Satoshi-Tb/go_todo_app/entity"
	"github.com/Satoshi-Tb/go_todo_app/testutil"
	"github.com/Satoshi-Tb/go_todo_app/testutil/fixture"
	"github.com/google/go-cmp/cmp"
)

func TestRepository_ListTasks(t *testing.T) {
	ctx := context.Background()

	// entity.Taskを作成する他のテストケースと混ざるとテストがフェイルする。
	// そのため、トランザクションをはることでこのテストケースの中だけのテーブル状態にする。
	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)

	//テスト後にロールバックする
	t.Cleanup(func() { _ = tx.Rollback() })
	if err != nil {
		t.Fatal(err)
	}

	wantUserID, wants := prepareTasks(ctx, t, tx)

	sut := &Repository{}

	gots, err := sut.ListTasks(ctx, tx, wantUserID)
	if err != nil {
		t.Fatalf("unexected error: %v", err)
	}
	if d := cmp.Diff(gots, wants); len(d) != 0 {
		t.Errorf("differs: (-got +want)\n%s", d)
	}

}

func prepareUser(ctx context.Context, t *testing.T, db Execer) entity.UserID {
	t.Helper()
	u := fixture.User(nil)
	result, err := db.ExecContext(ctx, "INSERT INTO user (name, password, role, created, modified) VALUES (?, ?, ?, ?, ?);", u.Name, u.Password, u.Role, u.Created, u.Modified)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("got user_id: %v", err)
	}
	t.Logf("use create succeeded. got user_id: %d", id)
	return entity.UserID(id)
}

func prepareTasks(ctx context.Context, t *testing.T, con Execer) (entity.UserID, entity.Tasks) {
	t.Helper()

	userID := prepareUser(ctx, t, con)
	otherUserID := prepareUser(ctx, t, con)

	c := clock.FixedClocker{}
	wants := entity.Tasks{
		{
			UserID: userID,
			Title:  "want task 1", Status: "todo",
			Created: c.Now(), Modified: c.Now(),
		},
		{
			UserID: userID,
			Title:  "want task 2", Status: "done",
			Created: c.Now(), Modified: c.Now(),
		},
	}
	tasks := entity.Tasks{
		wants[0],
		{
			UserID: otherUserID,
			Title:  "not want task", Status: "todo",
			Created: c.Now(), Modified: c.Now(),
		},
		wants[1],
	}

	//MySQLで複数行インサートする
	//LastInsertIdで返ってくるのは、最初にインサートした行
	result, err := con.ExecContext(ctx,
		`INSERT INTO task (user_id, title, status, created, modified)
			VALUES
			    (?, ?, ?, ?, ?),
			    (?, ?, ?, ?, ?),
			    (?, ?, ?, ?, ?);`,
		tasks[0].UserID, tasks[0].Title, tasks[0].Status, tasks[0].Created, tasks[0].Modified,
		tasks[1].UserID, tasks[1].Title, tasks[1].Status, tasks[1].Created, tasks[1].Modified,
		tasks[2].UserID, tasks[2].Title, tasks[2].Status, tasks[2].Created, tasks[2].Modified,
	)
	if err != nil {
		t.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	tasks[0].ID = entity.TaskID(id)
	tasks[1].ID = entity.TaskID(id + 1)
	tasks[2].ID = entity.TaskID(id + 2)

	return userID, wants
}

func TestRepository_DelTask(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)

	//テスト後にロールバックする
	t.Cleanup(func() { _ = tx.Rollback() })
	if err != nil {
		t.Fatal(err)
	}

	_, wants := prepareTasks(ctx, t, tx)

	sut := &Repository{}
	cnt, err := sut.DelTask(ctx, tx, wants[0])
	if err != nil {
		t.Fatalf("unexected error: %v", err)
	}

	if cnt != 1 {
		t.Errorf("deleted count not 1: %d", cnt)
	}
}

func TestRepository_UpdateTask(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)

	//テスト後にロールバックする
	t.Cleanup(func() { _ = tx.Rollback() })
	if err != nil {
		t.Fatal(err)
	}

	wantUserID, wants := prepareTasks(ctx, t, tx)
	t.Logf("UserID %d", wantUserID)
	wants[0].Title = "new title"
	wants[0].Status = entity.TaskStatusDone
	wants[1].Status = entity.TaskStatusDoing

	sut := &Repository{
		Clocker: clock.FixedClocker{},
	}

	for i, w := range wants {
		t.Logf("[%d] update item: %+v", i, w)
		cnt, err := sut.UpdateTask(ctx, tx, w)
		if err != nil {
			t.Fatalf("[%d] unexected error: %v", i, err)
		}

		if cnt != 1 {
			t.Errorf("[%d] updated count not 1, but %d", i, cnt)
		}

	}
}
