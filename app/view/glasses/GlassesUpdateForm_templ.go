// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.707
package glasses

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"fmt"
	"github.com/FACorreiaa/glasses-management-platform/app/models"
)

func GlassesUpdateForm(form models.GlassesForm, id string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<section class=\"w-full bg-white\"><div class=\"mx-auto max-w-7xl\"><h2 class=\"px-4 mb-4 text-5xl font-bold text-gray-900 xl:text-6xl mb-10\">Update</h2><div class=\"flex flex-col lg:flex-row\"><form method=\"put\" hx-target=\"#success-message\"><input type=\"hidden\" name=\"_method\" value=\"PUT\"><div class=\"flex flex-row mb-10\"><div class=\"w-50 md:w-1/2 px-4 mb-8\"><label class=\"font-medium text-slate-900\">Reference</label> <input type=\"text\" class=\"block w-50 px-4 py-4 mt-1 text-xl placeholder-gray-400 bg-gray-200 rounded-lg focus:outline-none focus:ring-4 focus:ring-blue-600 focus:ring-opacity-50\" data-primary=\"blue-600\" data-rounded=\"rounded-lg\" placeholder=\"Insert reference\" name=\"reference\" autocomplete=\"reference\" id=\"reference\" required value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(form.Values["Reference"])
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `app/view/glasses/GlassesUpdateForm.templ`, Line: 28, Col: 40}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"></div><div class=\"w-50 md:w-1/2 px-4 mb-8\"><label class=\"font-medium text-slate-900\">Brand</label> <input type=\"text\" class=\"block w-50 px-4 py-4 mt-1 text-xl placeholder-gray-400 bg-gray-200 rounded-lg focus:outline-none focus:ring-4 focus:ring-blue-600 focus:ring-opacity-50\" data-primary=\"blue-600\" data-rounded=\"rounded-lg\" placeholder=\"Insert brand\" name=\"brand\" autocomplete=\"brand\" id=\"brand\" required value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(form.Values["Brand"])
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `app/view/glasses/GlassesUpdateForm.templ`, Line: 43, Col: 36}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"></div></div><div class=\"flex flex-row mb-10\"><div class=\"w-50 md:w-1/2 px-4 mb-8\"><label class=\"font-medium text-slate-900\">Left eye strength</label> <input type=\"number\" class=\"block w-50 px-4 py-4 mt-1 text-xl placeholder-gray-400 bg-gray-200 rounded-lg focus:outline-none focus:ring-4 focus:ring-blue-600 focus:ring-opacity-50\" data-primary=\"blue-600\" data-rounded=\"rounded-lg\" placeholder=\"Insert graduation\" name=\"left_eye_strength\" autocomplete=\"left-eye\" id=\"left-eye\" min=\"-99\" max=\"99\" step=\"0.1\" required value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 string
		templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(form.Values["LeftEye"])
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `app/view/glasses/GlassesUpdateForm.templ`, Line: 63, Col: 38}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"></div><div class=\"w-50 md:w-1/2 px-4 mb-8\"><label class=\"font-medium text-slate-900\">Right eye strength</label> <input type=\"number\" class=\"block w-50 px-4 py-4 mt-1 text-xl placeholder-gray-400 bg-gray-200 rounded-lg focus:outline-none focus:ring-4 focus:ring-blue-600 focus:ring-opacity-50\" data-primary=\"blue-600\" data-rounded=\"rounded-lg\" placeholder=\"Insert graduation\" name=\"right_eye_strength\" autocomplete=\"right-eye\" id=\"right-eye\" min=\"-99\" max=\"99\" step=\"0.1\" required value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var5 string
		templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(form.Values["RightEye"])
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `app/view/glasses/GlassesUpdateForm.templ`, Line: 81, Col: 39}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"></div></div><div class=\"flex flex-row mb-10\"><div class=\"w-50 md:w-1/2 px-4 mb-8\"><label class=\"font-medium text-slate-900\">Color</label> <input type=\"text\" class=\"block w-50 px-4 py-4 mt-1 text-xl placeholder-gray-400 bg-gray-200 rounded-lg focus:outline-none focus:ring-4 focus:ring-blue-600 focus:ring-opacity-50\" data-primary=\"blue-600\" data-rounded=\"rounded-lg\" placeholder=\"Insert color\" name=\"color\" autocomplete=\"color\" id=\"color\" required value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var6 string
		templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(form.Values["Color"])
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `app/view/glasses/GlassesUpdateForm.templ`, Line: 98, Col: 36}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"></div><div class=\"w-50 md:w-1/2 px-4 mb-8\"><label class=\"font-medium text-slate-900\">Type</label> <select class=\"select block w-full max-w-xs mt-1 text-xl placeholder-gray-400 bg-gray-200 rounded-lg focus:outline-none focus:ring-4 focus:ring-blue-600 focus:ring-opacity-50\" data-primary=\"blue-600\" data-rounded=\"rounded-lg\" name=\"type\" id=\"type\" required><option value=\"adult\">Adult</option> <option value=\"children\">Children</option></select></div></div><div class=\"relative px-4 w-full mb-4\"><label class=\"font-medium text-slate-900\">Features</label> <textarea type=\"text\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var7 string
		templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(form.Values["Feature"])
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `app/view/glasses/GlassesUpdateForm.templ`, Line: 120, Col: 37}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" name=\"features\" value=\"features\" placeholder=\"Glasses features\" class=\"mt-1 flex relative z-20 peer w-full h-auto min-h-[80px] px-3 py-2 text-sm bg-white border-2 border-neutral-900 placeholder:text-neutral-500 focus:text-neutral-800 focus:border-neutral-900 focus:outline-none focus:ring-0 disabled:cursor-not-allowed disabled:opacity-50\"></textarea><div class=\"absolute inset-0 z-10 w-full h-full -m-1 duration-300 ease-out translate-x-2 translate-y-2 bg-black peer-focus:m-0 peer-focus:translate-x-0 peer-focus:translate-y-0\"></div></div><div class=\"flex flex-row\"><div class=\"w-50 md:w-1/2 px-4 mb-8\"><button hx-put=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var8 string
		templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(string(templ.URL(fmt.Sprintf("/glasses/%s/update", id))))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `app/view/glasses/GlassesUpdateForm.templ`, Line: 131, Col: 73}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" type=\"submit\" name=\"action\" class=\"btn btn-primary w-full max-w-xs\" value=\"update_and_redirect\">Update and Go to List</button></div></div><div id=\"success-message\" class=\"mt-4 text-success\"></div></form></div></div></section><script>\n        document.addEventListener(\"glassesUpdated\", function() {\n            const successMessage = document.getElementById(\"success-message\");\n            if (successMessage) {\n                setTimeout(() => {\n                    successMessage.style.display = \"none\";\n                }, 3000);\n            }\n        });\n    </script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}