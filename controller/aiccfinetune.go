package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-aicc-finetune/app"
)

func AddRouterForAICCFinetuneController(
	rg *gin.RouterGroup,
	fs app.FinetuneService,
) {
	ctl := AICCFinetuneController{fs: fs}

	rg.POST("/v1/aiccfinetune", ctl.Create)
	rg.DELETE("/v1/aiccfinetune/:id", ctl.Delete)
	rg.PUT("/v1/aiccfinetune/:id", ctl.Terminate)
	rg.GET("/v1/aiccfinetune/:id/log", ctl.GetLog)
	rg.GET("/v1/aiccfinetune/:id/result/:file", ctl.GetDownloadURL)

}

type AICCFinetuneController struct {
	baseController

	fs app.FinetuneService
}

// @Summary		Create
// @Description	create aicc finetune
// @Tags			AICC Finetune
// @Param			body	body	AICCFinetuneCreateRequest	true	"body of creating aicc finetune"
// @Accept			json
// @Success		201	{object}			app.AICCFinetuneInfoDTO
// @Failure		400	bad_request_body	can't	parse		request	body
// @Failure		401	bad_request_param	some	parameter	of		body	is	invalid
// @Failure		500	system_error		system	error
// @Router			/v1/aiccfinetune [post]
func (ctl *AICCFinetuneController) Create(ctx *gin.Context) {
	req := AICCFinetuneCreateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}
	cmd := new(app.AICCFinetuneCreateCmd)
	err := req.toCmd(cmd)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))
		return
	}

	v, err := ctl.fs.Create(cmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusCreated, newResponseData(v))
}

// @Summary		Create
// @Description	create aicc finetune
// @Tags			AICC Finetune
// @Param			body	body	AICCFinetuneCreateRequest	true	"body of creating aicc finetune"
// @Accept			json
// @Success		201	{object}			app.AICCFinetuneInfoDTO
// @Failure		400	bad_request_body	can't	parse		request	body
// @Failure		401	bad_request_param	some	parameter	of		body	is	invalid
// @Failure		500	system_error		system	error
// @Router			/v1/aiccfinetune [post]
func (ctl *AICCFinetuneController) Delete(ctx *gin.Context) {
	jobId := ctx.Param("id")
	if err := ctl.fs.Delete(jobId); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusNoContent, newResponseData("success"))
}

// @Summary		Terminate
// @Description	terminate aicc finetune job
// @Tags			AICC Finetune
// @Param			body	body	AICCFinetuneCreateRequest	true	"body of creating aicc finetune"
// @Accept			json
// @Success		201	{object}			app.AICCFinetuneInfoDTO
// @Failure		400	bad_request_body	can't	parse		request	body
// @Failure		401	bad_request_param	some	parameter	of		body	is	invalid
// @Failure		500	system_error		system	error
// @Router			/v1/aiccfinetune/{id} [put]
func (ctl *AICCFinetuneController) Terminate(ctx *gin.Context) {
	jobId := ctx.Param("id")
	if err := ctl.fs.Terminate(jobId); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData("success"))
}

// @Summary		GetLog
// @Description	get log url of aicc finetune for downloading
// @Tags			AICC Finetune
// @Param			id	path	string	true	"id of aicc finetune job"
// @Accept			json
// @Success		200	{object}		AICCFinetuneResultResp
// @Failure		500	system_error	system	error
// @Router			/v1/aiccfinetune/{id}/log [get]
func (ctl *AICCFinetuneController) GetLog(ctx *gin.Context) {
	v, err := ctl.fs.GetLogDownloadURL(ctx.Param("id"))
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(AICCFinetuneResultResp{v}))
}

// @Summary		GetDownloadURL
// @Description	get download url of aicc finetune result such as log or output.
// @Tags			AICC Finetune
// @Param			id		path	string	true	"id of finetune"
// @Param			file	path	string	true	"obs file path to download"
// @Accept			json
// @Success		200	{object}		AICCFinetuneResultResp
// @Failure		500	system_error	system	error
// @Router			/v1/aiccfinetune/{id}/result/{file} [get]
func (ctl *AICCFinetuneController) GetDownloadURL(ctx *gin.Context) {
	v, err := ctl.fs.GenFileDownloadURL(ctx.Param("file"))
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(AICCFinetuneResultResp{v}))
}
