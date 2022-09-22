package applet

import (
	"context"
	"mime/multipart"
	"os"

	"github.com/google/uuid"
	"github.com/iotexproject/Bumblebee/kit/sqlx/builder"
	"github.com/iotexproject/Bumblebee/kit/sqlx/datatypes"

	"github.com/iotexproject/w3bstream/pkg/modules/resource"
	"github.com/iotexproject/w3bstream/pkg/types"

	"github.com/iotexproject/w3bstream/pkg/errors/status"
	"github.com/iotexproject/w3bstream/pkg/models"
)

type CreateAppletReq struct {
	File *multipart.FileHeader `name:"file"`
	Info `name:"info"`
}

type Info struct {
	ProjectID  string              `json:"projectID"`
	AppletName string              `json:"appletName"`
	Config     models.AppletConfig `json:"config,omitempty"`
}

func CreateApplet(ctx context.Context, r *CreateAppletReq) (*models.Applet, error) {
	appletID := uuid.New().String()
	_, filename, err := resource.Upload(ctx, r.File, appletID)
	if err != nil {
		return nil, status.UploadFileFailed.StatusErr().WithDesc(err.Error())
	}

	d := types.MustDBExecutorFromContext(ctx)

	m := &models.Applet{
		RelProject: models.RelProject{ProjectID: r.ProjectID},
		RelApplet:  models.RelApplet{AppletID: appletID},
		AppletInfo: models.AppletInfo{Name: r.AppletName, Path: filename, Config: r.Config},
	}

	if err = m.Create(d); err != nil {
		defer os.RemoveAll(filename)
		return nil, err
	}

	return m, nil
}

type ListAppletReq struct {
	ProjectID string   `in:"path"  name:"projectID"`
	IDs       []string `in:"query" name:"id,omitempty"`
	AppletIDs []string `in:"query" name:"appletID,omitempty"`
	Names     []string `in:"query" name:"name,omitempty"`
	datatypes.Pager
}

func (r *ListAppletReq) Condition() builder.SqlCondition {
	var (
		m  = &models.Applet{}
		cs []builder.SqlCondition
	)
	if len(r.IDs) > 0 {
		cs = append(cs, m.ColID().In(r.IDs))
	}
	if len(r.AppletIDs) > 0 {
		cs = append(cs, m.ColAppletID().In(r.AppletIDs))
	}
	if len(r.Names) > 0 {
		cs = append(cs, m.ColName().In(r.Names))
	}
	return builder.And(cs...)
}

func (r *ListAppletReq) Additions() builder.Additions {
	m := &models.Applet{}
	return builder.Additions{
		builder.OrderBy(builder.DescOrder(m.ColCreatedAt())),
		r.Pager.Addition(),
	}
}

type ListAppletRsp struct {
	Data  []models.Applet `json:"data"`
	Hints int64           `json:"hints"`
}

func ListApplets(ctx context.Context, r *ListAppletReq) (*ListAppletRsp, error) {
	applet := &models.Applet{}

	d := types.MustDBExecutorFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)

	l.Start(ctx, "ListApplets")
	defer l.End()

	applets, err := applet.List(d, r.Condition(), r.Additions()...)
	if err != nil {
		l.Error(err)
		return nil, err
	}
	hints, err := applet.Count(d, r.Condition())
	if err != nil {
		l.Error(err)
		return nil, err
	}
	return &ListAppletRsp{applets, hints}, nil
}

type RemoveAppletReq struct {
	ProjectID  string `in:"path"  name:"projectID"`
	AppletID   string `in:"query" name:"appletID,omitempty"`
	AppletName string `in:"query" name:"appletName,omitempty"`
}

func RemoveApplet(ctx context.Context, r *RemoveAppletReq) error {
	// TODO
	return nil
}

type GetAppletReq struct {
	ProjectID string `in:"path" name:"projectID"`
	AppletID  string `in:"path" name:"appletID"`
}

type GetAppletRsp struct {
	models.Applet
	Instances []models.Instance `json:"instances"`
}

func GetAppletByID(ctx context.Context, appletID string) (*GetAppletRsp, error) {
	d := types.MustDBExecutorFromContext(ctx)
	m := &models.Applet{RelApplet: models.RelApplet{AppletID: appletID}}
	err := m.FetchByAppletID(d)
	return &GetAppletRsp{
		Applet: *m,
	}, err
}
