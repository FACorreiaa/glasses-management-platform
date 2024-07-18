// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package glasses

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"github.com/FACorreiaa/glasses-management-platform/app/models"
)

func GlassesRegisterForm(form models.GlassesForm) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<section class=\"w-full bg-white py-8\"><div class=\"mx-auto max-w-7xl px-4 sm:px-6 lg:px-8\"><h2 class=\"mb-6 text-3xl font-extrabold text-gray-900 xl:text-4xl\">Insert Glasses</h2><div class=\"flex flex-col lg:flex-row\"><form method=\"post\" hx-target=\"#success-message\" class=\"w-full space-y-4\"><div class=\"flex flex-wrap -mx-4\"><div class=\"w-full md:w-1/2 px-4 mb-4\"><label class=\"block font-medium text-gray-900 mb-1\" for=\"reference\">Reference</label> <input type=\"text\" class=\"block w-full px-4 py-3 mt-1 text-base placeholder-gray-400 bg-gray-100 rounded-lg focus:outline-none \" placeholder=\"Insert reference\" name=\"reference\" autocomplete=\"reference\" id=\"reference\" required value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(form.Values["Reference"])
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `app/view/glasses/GlassesRegisterForm.templ`, Line: 24, Col: 40}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"></div><div class=\"w-full md:w-1/2 px-4 mb-4\"><label class=\"block font-medium text-gray-900 mb-1\" for=\"brand\">Brand</label> <input type=\"text\" class=\"block w-full px-4 py-3 mt-1 text-base placeholder-gray-400 bg-gray-100 rounded-lg focus:outline-none \" placeholder=\"Insert brand\" name=\"brand\" autocomplete=\"brand\" id=\"brand\" required value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(form.Values["Brand"])
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `app/view/glasses/GlassesRegisterForm.templ`, Line: 37, Col: 36}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"></div></div><div class=\"flex flex-wrap -mx-4\"><div class=\"w-full md:w-1/2 px-4 mb-4\"><label class=\"block font-medium text-gray-900 mb-1\" for=\"left-eye\">Left Eye Strength</label> <input type=\"number\" class=\"block w-full px-4 py-3 mt-1 text-base placeholder-gray-400 bg-gray-100 rounded-lg focus:outline-none \" placeholder=\"Insert graduation\" name=\"left_eye_strength\" autocomplete=\"left-eye\" id=\"left-eye\" min=\"-99\" max=\"99\" step=\"0.1\" required value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 string
		templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(form.Values["LeftEye"])
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `app/view/glasses/GlassesRegisterForm.templ`, Line: 55, Col: 38}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"></div><div class=\"w-full md:w-1/2 px-4 mb-4\"><label class=\"block font-medium text-gray-900 mb-1\" for=\"right-eye\">Right Eye Strength</label> <input type=\"number\" class=\"block w-full px-4 py-3 mt-1 text-base placeholder-gray-400 bg-gray-100 rounded-lg focus:outline-none \" placeholder=\"Insert graduation\" name=\"right_eye_strength\" autocomplete=\"right-eye\" id=\"right-eye\" min=\"-99\" max=\"99\" step=\"0.1\" required value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var5 string
		templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(form.Values["RightEye"])
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `app/view/glasses/GlassesRegisterForm.templ`, Line: 71, Col: 39}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"></div></div><div class=\"flex flex-wrap -mx-4\"><div class=\"w-full md:w-1/2 px-4 mb-4\"><label class=\"block font-medium text-gray-900 mb-1\" for=\"color\">Color</label> <input type=\"text\" class=\"block w-full px-4 py-3 mt-1 text-base placeholder-gray-400 bg-gray-100 rounded-lg focus:outline-none \" placeholder=\"Insert color\" name=\"color\" autocomplete=\"color\" id=\"color\" required value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var6 string
		templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(form.Values["Color"])
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `app/view/glasses/GlassesRegisterForm.templ`, Line: 86, Col: 36}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"></div><div class=\"w-full md:w-1/2 px-4 mb-4\"><label class=\"block font-medium text-gray-900 mb-1\" for=\"type\">Type</label> <select class=\"block w-full px-4 py-3 mt-1 text-base placeholder-gray-400 bg-gray-100 rounded-lg focus:outline-none \" name=\"type\" id=\"type\" required value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var7 string
		templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(form.Values["Type"])
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `app/view/glasses/GlassesRegisterForm.templ`, Line: 96, Col: 35}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><option disabled selected>Select type</option> <option value=\"adult\">Adult</option> <option value=\"children\">Children</option></select></div></div><div class=\"relative w-full mb-4\"><label class=\"block font-medium text-gray-900 mb-1\" for=\"features\">Features</label> <textarea class=\"block w-full px-4 py-3 mt-1 text-base placeholder-gray-400 bg-gray-100 rounded-lg focus:outline-none \" name=\"features\" placeholder=\"Glasses features\" id=\"features\" required value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var8 string
		templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(form.Values["Feature"])
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `app/view/glasses/GlassesRegisterForm.templ`, Line: 112, Col: 37}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"></textarea></div><div class=\"flex flex-wrap -mx-4\"><div class=\"w-full md:w-1/2 px-4 mb-4\"><button type=\"submit\" class=\"btn w-full px-4 py-3 bg-blue-500 text-white text-base font-medium rounded-lg hover:bg-blue-600 focus:outline-none \" name=\"action\" value=\"insert_more\">Insert and Add Another</button></div><div class=\"w-full md:w-1/2 px-4 mb-4\"><button type=\"submit\" class=\"btn w-full px-4 py-3 bg-blue-500 text-white text-base font-medium rounded-lg hover:bg-blue-600 focus:outline-none \" name=\"action\" value=\"insert_and_redirect\">Insert and Go to List</button></div></div><div id=\"success-message\" class=\"mt-4 text-green-500\"></div></form></div></div></section><script>\n\t\tdocument.addEventListener(\"submit\", function(event) {\n\t\t\tconst form = event.target;\n\t\t\tif (form.getAttribute('method') === 'post') {\n\t\t\t\tconst successMessage = document.getElementById(\"success-message\");\n\n\t\t\t\t// Display the success message\n\t\t\t\tif (successMessage) {\n\t\t\t\t\tsuccessMessage.textContent = \"Item inserted successfully!\";\n\t\t\t\t\tsuccessMessage.style.display = \"block\";\n\n\t\t\t\t\t// Hide the success message after 3 seconds\n\t\t\t\t\tsetTimeout(() => {\n\t\t\t\t\t\tsuccessMessage.style.display = \"none\";\n\t\t\t\t\t}, 3000);\n\n\t\t\t\t\t// Trigger custom event to handle additional logic if necessary\n\t\t\t\t\tconst glassesAddedEvent = new Event(\"glassesAdded\");\n\t\t\t\t\tdocument.dispatchEvent(glassesAddedEvent);\n\t\t\t\t}\n\t\t\t}\n\t\t});\n\t</script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
