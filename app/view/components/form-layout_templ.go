// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.707
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func FormLayoutTemplate() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<section class=\"w-full bg-white\"><div class=\"mx-auto max-w-7xl\"><div class=\"flex flex-col lg:flex-row\"><div class=\"relative w-full bg-cover lg:w-6/12 xl:w-7/12 bg-gradient-to-r from-white via-white to-gray-100\"><div class=\"relative flex flex-col items-center justify-center w-full h-full px-10 my-20 lg:px-16 lg:my-0\"><div class=\"flex flex-col items-start space-y-8 tracking-tight lg:max-w-3xl\"><div class=\"relative\"><p class=\"mb-2 font-medium text-gray-700 uppercase\">Work smarter</p><h2 class=\"text-5xl font-bold text-gray-900 xl:text-6xl\">Features to help you work smarter</h2></div><p class=\"text-2xl text-gray-700\">We've created a simple formula to follow in order to gain more out of your business and your application.</p><a href=\"#_\" class=\"inline-block px-8 py-5 text-xl font-medium text-center text-white transition duration-200 bg-blue-600 rounded-lg hover:bg-blue-700 ease\" data-primary=\"blue-600\" data-rounded=\"rounded-lg\">Get Started Today</a></div></div></div><div class=\"w-full bg-white lg:w-6/12 xl:w-5/12\"><div class=\"flex flex-col items-start justify-start w-full h-full p-10 lg:p-16 xl:p-24\"><h4 class=\"w-full text-3xl font-bold\">Signup</h4><p class=\"text-lg text-gray-500\">or, if you have an account you can <a href=\"#_\" class=\"text-blue-600 underline\" data-primary=\"blue-600\">sign in</a></p><div class=\"relative w-full mt-10 space-y-8\"><div class=\"relative\"><label class=\"font-medium text-gray-900\">Name</label> <input type=\"text\" class=\"block w-full px-4 py-4 mt-2 text-xl placeholder-gray-400 bg-gray-200 rounded-lg focus:outline-none focus:ring-4 focus:ring-blue-600 focus:ring-opacity-50\" data-primary=\"blue-600\" data-rounded=\"rounded-lg\" placeholder=\"Enter Your Name\"></div><div class=\"relative\"><label class=\"font-medium text-gray-900\">Email</label> <input type=\"text\" class=\"block w-full px-4 py-4 mt-2 text-xl placeholder-gray-400 bg-gray-200 rounded-lg focus:outline-none focus:ring-4 focus:ring-blue-600 focus:ring-opacity-50\" data-primary=\"blue-600\" data-rounded=\"rounded-lg\" placeholder=\"Enter Your Email Address\"></div><div class=\"relative\"><label class=\"font-medium text-gray-900\">Password</label> <input type=\"password\" class=\"block w-full px-4 py-4 mt-2 text-xl placeholder-gray-400 bg-gray-200 rounded-lg focus:outline-none focus:ring-4 focus:ring-blue-600 focus:ring-opacity-50\" data-primary=\"blue-600\" data-rounded=\"rounded-lg\" placeholder=\"Password\"></div><div class=\"relative\"><a href=\"#_\" class=\"inline-block w-full px-5 py-4 text-lg font-medium text-center text-white transition duration-200 bg-blue-600 rounded-lg hover:bg-blue-700 ease\" data-primary=\"blue-600\" data-rounded=\"rounded-lg\">Create Account</a> <a href=\"#_\" class=\"inline-block w-full px-5 py-4 mt-3 text-lg font-bold text-center text-gray-900 transition duration-200 bg-white border border-gray-300 rounded-lg hover:bg-gray-100 ease\" data-rounded=\"rounded-lg\">Sign up with Google</a></div></div></div></div></div></div></section>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
