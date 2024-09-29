package pages

import (
    "os"
	"io"
	"net/http"
	"nihility/dbase"
	"nihility/logger"
	"nihility/logic"
	"strconv"
    "github.com/gorilla/sessions"
)

var log = logger.Logger {
    Color: logger.Colors.Red,
    Pretext: "pages",
}

type sessioner struct {
    Auth logic.Auth
    Main string
    Path string
    Dto any
}
//FIXME: Handle fully separately in every function/session!!
//var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
var store = sessions.NewCookieStore([]byte(os.Getenv("NYANTAN_SESSION_KEY")))
var session sessioner

var artifact_path string = "artifacts/"
var html_path string = "html/"
var base_template_path string = html_path + "base.html"

func authenticate(r *http.Request) {
    // TODO: Add request aut header
    real_session, _ := store.Get(r, "uname")
    //log.Println(real_session)
    uname, _ := real_session.Values["uname"].(string)

    session.Auth.Username = uname
    logic.Authenticate(&session.Auth)
}

func Base_auth_and_render(w http.ResponseWriter, r *http.Request, path string) (string, string) {
    session.Path = r.URL.Path
    authenticate(r)
    return read_artifact(path, w.Header())
}

// =====================================================================================================================
// Basic functios

func Root(w http.ResponseWriter, r *http.Request) {
    if "/" == r.URL.Path {
        fil, _ := Base_auth_and_render(w, r, "index.html")
        Render(w, fil, nil)
    } else {
        Unexpected(w, r)
    }
}

func Login(w http.ResponseWriter, r *http.Request) {
    fil, _ := Base_auth_and_render(w, r, "login.html")
    uname := r.FormValue("form[userName]")
    upass := r.FormValue("form[userPass]")

    if "" != uname {
        if logic.Auth_login(uname, upass).Id != "" {
            // FIXME: Store auth headers in database with associated user
            rsess, _ := store.New(r, "uname")
            log.Printf("session is: %s\n", rsess.ID)
            rsess.Values["uname"] = uname
            rsess.Save(r, w)
            session.Auth.Username = uname
        } else {
            session.Auth.Error = "Auth Error"
        }
    } else {
        session.Auth.Username = ""
        session.Auth.Error = ""
    }

    if "" == session.Auth.Username {
        Render(w, fil, nil)
    } else {
        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}

func Register(w http.ResponseWriter, r *http.Request) {
    fil, _ := Base_auth_and_render(w, r, "login.html")
    uname := r.FormValue("form[userName]")
    upass := r.FormValue("form[userPass]")

    if "" != uname {
        if logic.Auth_register(uname, upass) {
            session.Auth.Username = uname
        } else {
            session.Auth.Error = "Cannot Register"
        }
    } else {
        session.Auth.Username = ""
    }
        session.Auth.Error = ""

    if "" == session.Auth.Username {
        Render(w, fil, nil)
    } else {
        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}

func Translate(w http.ResponseWriter, r *http.Request) {
    fil, _ := Base_auth_and_render(w, r, "translate.html")
    translations, err := logic.List_translations(session.Auth)
    if err != nil {
        log.Println(err)
    }
    Render(w, fil, translations)
}

func Logout(w http.ResponseWriter, r *http.Request) {
    logic.Authenticate(&session.Auth)
    rsess, _ := store.Get(r, "uname")
    rsess.Options.MaxAge = -1
    rsess.Save(r, w)
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Unexpected(w http.ResponseWriter, r *http.Request) {
    fil, typ := Base_auth_and_render(w, r, r.URL.Path)

    if "text" == typ {
        Render(w, fil, nil)
    } else {
        io.WriteString(w, fil)
    }
}

// =====================================================================================================================
// "Smart" functios

func base_error_render(w http.ResponseWriter, r *http.Request) {
    fil, _ := Base_auth_and_render(w, r, "not_found.html")
    Render(w, fil, nil)
}

func Translation(w http.ResponseWriter, r *http.Request) {
    selected, err := dbase.Select_translation(r.PathValue("id"))
    if nil != err {
        base_error_render(w, r)
        return
    }

    fil, _ := Base_auth_and_render(w, r, "trans.html")
    pre_rendered := Pre_render(fil, selected)
    Render(w, pre_rendered, nil)
}

func Editor_list(w http.ResponseWriter, r *http.Request) {
    fil, _ := Base_auth_and_render(w, r, "edit-list.html")
    id := r.PathValue("id")

    edits, err := logic.List_edits(id)
    if nil != err {
        base_error_render(w, r)
        return
    }

    epl := logic.Edit_page_list{
        TransId: id,
        Title: id,
        Link: logic.Generate_translation_link(id),
        PageCount: len(edits),
        Edits: edits,
    }

    pre_rendered := Pre_render(fil, epl)
    Render(w, pre_rendered, nil)
}

func Editor(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    page := r.PathValue("page")

    selected, err := dbase.Select_translation(id)
    if nil != err {
        base_error_render(w, r)
        return
    }

    if !logic.User_in_fandom(session.Auth, selected.Fandom) {
        base_error_render(w, r)
        return
    }

    log.Println(id, page)
    if page == "" {
        Editor_list(w, r)
        return
    }

    page_index, err := strconv.Atoi(page)
    if nil != err {
        base_error_render(w, r)
        return
    }

    edits, err := logic.Select_edit(id, page_index)
    if nil != err {
        base_error_render(w, r)
        return
    }
    edit_list := logic.Edit_list{
        TransId:    selected.Id.String(),
        // FIXME: sould be setted with prefixes and paths
        Title:      selected.Title,
        Link:       selected.Link,
        Image:      logic.Generate_translation_image_path_original(id, page_index),
        Page:       page_index,
        PageCount:  selected.Pages,
        Edits:      edits,
    }

    fil, _ := Base_auth_and_render(w, r, "editor.html")
    pre_rendered := Pre_render(fil, edit_list)
    Render(w, pre_rendered, nil)
}
