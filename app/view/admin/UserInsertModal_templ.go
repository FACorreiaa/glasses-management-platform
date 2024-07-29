// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package admin

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"github.com/FACorreiaa/glasses-management-platform/app/models"
	"github.com/FACorreiaa/glasses-management-platform/app/view/components"
)

func UserInsertModal(form models.RegisterFormValues) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div x-data=\"{ modalOpen: false }\" @keydown.escape.window=\"modalOpen = false\" :class=\"{ &#39;z-40&#39;: modalOpen }\" class=\"relative w-auto h-auto\"><button @click=\"modalOpen=true\" class=\"btn btn-xs btn-primary inline-flex items-center justify-center text-sm font-medium tracking-wide transition-colors duration-200 rounded-md hover:bg-neutral-900 focus:ring-2 focus:ring-offset-2 focus:ring-neutral-900 focus:shadow-outline focus:outline-none\">Insert collaborator</button><template x-teleport=\"body\"><div x-show=\"modalOpen\" class=\"fixed top-0 left-0 z-[99] flex items-center justify-center w-screen h-screen\" x-cloak><div x-show=\"modalOpen\" x-transition:enter=\"ease-out duration-300\" x-transition:enter-start=\"opacity-0\" x-transition:enter-end=\"opacity-100\" x-transition:leave=\"ease-in duration-300\" x-transition:leave-start=\"opacity-100\" x-transition:leave-end=\"opacity-0\" @click=\"modalOpen=false\" class=\"absolute inset-0 w-full h-full bg-gray-900 bg-opacity-50 backdrop-blur-sm\"></div><div x-show=\"modalOpen\" x-trap.inert.noscroll=\"modalOpen\" x-transition:enter=\"ease-out duration-300\" x-transition:enter-start=\"opacity-0 scale-90\" x-transition:enter-end=\"opacity-100 scale-100\" x-transition:leave=\"ease-in duration-200\" x-transition:leave-start=\"opacity-100 scale-100\" x-transition:leave-end=\"opacity-0 scale-90\" class=\"relative w-full py-6 bg-white shadow-md px-7 bg-opacity-90 drop-shadow-md backdrop-blur-sm sm:max-w-lg sm:rounded-lg\"><div class=\"flex items-center justify-between pb-3\"><button @click=\"modalOpen=false\" class=\"absolute top-0 right-0 flex items-center justify-center w-8 h-8 mt-5 mr-5 text-gray-600 rounded-full hover:text-gray-800 hover:bg-gray-50\"><svg class=\"w-5 h-5\" xmlns=\"http://www.w3.org/2000/svg\" fill=\"none\" viewBox=\"0 0 24 24\" stroke-width=\"1.5\" stroke=\"currentColor\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M6 18L18 6M6 6l12 12\"></path></svg></button></div><div class=\"relative w-auto pb-8\"><section class=\"w-full bg-white\"><div class=\"mx-auto max-w-7xl\"><h2 class=\"mb-8 text-5xl font-bold text-gray-900 xl:text-6xl mb-10\">Register user</h2><div class=\"flex flex-col lg:flex-row\"><form method=\"post\" class=\"w-full space-y-4\"><fieldset><label for=\"username\" class=\"block text-sm font-medium text-gray-900 dark:text-white\">Username</label> <input class=\"block w-full px-4 py-4 mt-2 text-xl placeholder-gray-400 bg-gray-200 rounded-lg focus:outline-none \" data-primary=\"blue-600\" data-rounded=\"rounded-lg\" type=\"text\" placeholder=\"Username\" required name=\"username\" autocomplete=\"username\" id=\"username\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(form.Values["username"])
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `app/view/admin/UserInsertModal.templ`, Line: 63, Col: 43}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if form.FieldErrors["username"] != "" {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<p class=\"text-error\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(form.FieldErrors["username"])
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `app/view/admin/UserInsertModal.templ`, Line: 66, Col: 64}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</fieldset><fieldset class=\"max-w-lg mx-auto mt-2\"><label for=\"email\" class=\"block text-sm font-medium text-gray-900 dark:text-white\">Email</label> <input class=\"block w-full px-4 py-4 mt-2 text-xl placeholder-gray-400 bg-gray-200 rounded-lg focus:outline-none \" data-primary=\"blue-600\" data-rounded=\"rounded-lg\" type=\"email\" placeholder=\"Email\" required name=\"email\" id=\"email\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 string
		templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(form.Values["email"])
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `app/view/admin/UserInsertModal.templ`, Line: 80, Col: 40}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"></fieldset><fieldset class=\"max-w-lg mx-auto mt-2\"><label for=\"passoword\" class=\"block text-sm font-medium text-gray-900 dark:text-white\">Password</label> <input class=\"block w-full px-4 py-4 mt-2 text-xl placeholder-gray-400 bg-gray-200 rounded-lg focus:outline-none \" data-primary=\"blue-600\" data-rounded=\"rounded-lg\" type=\"password\" placeholder=\"Password\" required name=\"password\" autocomplete=\"new-password\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var5 string
		templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(form.Values["password"])
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `app/view/admin/UserInsertModal.templ`, Line: 94, Col: 43}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if form.FieldErrors["password"] != "" {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<p class=\"text-error\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var6 string
			templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(form.FieldErrors["password"])
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `app/view/admin/UserInsertModal.templ`, Line: 97, Col: 64}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</fieldset><fieldset class=\"max-w-lg mx-auto mt-2\"><label for=\"password\" class=\"block text-sm font-medium text-gray-900 dark:text-white\">Confirm Password</label> <input class=\"block w-full px-4 py-4 mt-2 text-xl placeholder-gray-400 bg-gray-200 rounded-lg focus:outline-none\" data-primary=\"blue-600\" data-rounded=\"rounded-lg\" type=\"password\" placeholder=\"Confirm Password\" required name=\"password_confirm\" autocomplete=\"new-password\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var7 string
		templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(form.Values["password_confirm"])
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `app/view/admin/UserInsertModal.templ`, Line: 111, Col: 51}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if form.FieldErrors["password_confirm"] != "" {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<p class=\"text-error\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var8 string
			templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(form.FieldErrors["password_confirm"])
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `app/view/admin/UserInsertModal.templ`, Line: 114, Col: 72}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</fieldset><div class=\"flex items-center justify-center max-w-lg mx-auto mt-2\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.ButtonCancelComponent().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<button hx-post=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var9 string
		templ_7745c5c3_Var9, templ_7745c5c3_Err = templ.JoinStringErrs("/collaborators/register")
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `app/view/admin/UserInsertModal.templ`, Line: 120, Col: 47}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var9))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" type=\"submit\" class=\"btn btn-xs btn-primary inline-flex items-center justify-center text-sm font-medium tracking-wide transition-colors duration-200 rounded-md hover:bg-neutral-900 focus:ring-2 focus:ring-offset-2 focus:ring-neutral-900 focus:shadow-outline focus:outline-none mr-2\" name=\"action\" value=\"submit\">Insert</button></div></form></div></div></section></div></div></div></template></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
