package main

import "path/filepath"

var fileKind = map[string]string{
		".txt":   "Plain Text",
	".doc":   "Microsoft Word Document",
	".docx":  "Microsoft Word Document",
	".log":   "Log File",
	".msg":   "Outlook Mail Message",
	".odt":   "OpenDocument Text Document",
	".pages": "Pages Document",
	".rtf":   "Rich Text Format",
	".tex":   "LaTeX Source Document",
	".wpd":   "WordPerfect Document",
	".wps":   "Microsoft Works Word Processor Document",
	".csv":   "Comma Separated Values File",
	".dat":   "Data File",
	".ged":   "GEDCOM Genealogy Data File",
	".key":   "Keynote Presentation",
	".keychain": "Mac OS X Keychain File",
	".pps":   "PowerPoint Slide Show",
	".ppt":   "PowerPoint Presentation",
	".pptx":  "PowerPoint Presentation",
	".sdf":   "Standard Data File",
	".tar":   "Consolidated Unix File Archive",
	".vcf":   "vCard File",
	".xml":   "XML File",
	".aif":   "Audio Interchange File Format",
	".iff":   "Interchange File Format",
	".m3u":   "Media Playlist File",
	".m4a":   "MPEG-4 Audio File",
	".mid":   "MIDI File",
	".mp3":   "MP3 Audio File",
	".mpa":   "MPEG-2 Audio File",
	".wav":   "WAVE Audio File",
	".wma":   "Windows Media Audio File",
	".3g2":   "3GPP2 Multimedia File",
	".3gp":   "3GPP Multimedia File",
	".asf":   "Advanced Systems Format File",
	".avi":   "Audio Video Interleave File",
	".flv":   "Flash Video File",
	".m4v":   "iTunes Video File",
	".mov":   "Apple QuickTime Movie",
	".mp4":   "MPEG-4 Video File",
	".mpg":   "MPEG Video File",
	".rm":    "Real Media File",
	".swf":   "Shockwave Flash Movie",
	".vob":   "DVD Video Object File",
	".wmv":   "Windows Media Video File",
	".bmp":   "Bitmap Image",
	".dds":   "DirectDraw Surface",
	".gif":   "Graphical Interchange Format File",
	".jpg":   "JPEG Image",
	".png":   "PNG Image",
	".psd":   "Adobe Photoshop Document",
	".pspimage": "PaintShop Pro Image",
	".tga":   "Targa Graphic",
	".thm":   "Thumbnail Image File",
	".tif":   "Tagged Image File",
	".tiff":  "Tagged Image File Format",
	".yuv":   "YUV Encoded Image File",
	".ai":    "Adobe Illustrator File",
	".eps":   "Encapsulated PostScript File",
	".pdf":   "PDF Document",
	".pict":  "Picture File",
	".svg":   "Scalable Vector Graphics File",
	".indd":  "Adobe InDesign Document",
	".pct":   "Picture File",
	".xlr":   "Works Spreadsheet",
	".xls":   "Excel Spreadsheet",
	".xlsx":  "Excel Spreadsheet",
	".css":   "CSS Style Sheet",
	".htm":   "HTML File",
	".html":  "HTML File",
	".js":    "JavaScript File",
	".jsp":   "Java Server Page",
	".php":   "PHP Source Code File",
	".rss":   "Rich Site Summary",
	".xhtml": "XHTML File",
	".c":     "C/C++ Source Code File",
	".class": "Java Class File",
	".cpp":   "C++ Source Code File",
	".cs":    "C# Source Code File",
	".dtd":   "Document Type Definition File",
	".fla":   "Adobe Animate Animation",
	".java":  "Java Source Code File",
	".lua":   "Lua Source File",
	".m":     "Objective-C Implementation File",
	".pl":    "Perl Script",
	".py":    "Python Script",
	".sh":    "Bash Shell Script",
	".sln":   "Visual Studio Solution File",
	".swift": "Swift Source Code File",
	".vcxproj":"Visual C++ Project File",
	".xcodeproj": "Xcode Project",
	".bak":   "Backup File",
	".cab":   "Windows Cabinet File",
	".cfg":   "Configuration File",
	".cpl":   "Windows Control Panel Item",
	".cur":   "Windows Cursor",
	".dll":   "Dynamic Link Library",
	".dmp":   "Dump File",
	".drv":   "Device Driver",
	".icns":  "macOS Icon Resource File",
	".ico":   "Icon File",
	".ini":   "Initialization File",
	".lnk":   "Windows Shortcut",
	".msi":   "Windows Installer Package",
	".sys":   "Windows System File",
	".tmp":   "Temporary File",
	".3dm":   "3D Model File",
	".3ds":   "3D Studio Scene",
	".max":   "3ds Max Scene File",
	".obj":   "Wavefront 3D Object File",
	".dwg":   "AutoCAD Drawing Database File",
	".dxf":   "Drawing Exchange Format File",
	".gpx":   "GPS Exchange File",
	".kml":   "Keyhole Markup Language File",
	".kmz":   "Google Earth Placemark File",
	".webloc": "Website Shortcut",
}

func getKind(path string) string {
	ext := filepath.Ext(path)
	kind, ok := fileKind[ext]
	if !ok {
		return "Unknown"  // Default kind
	}
	return kind
}
