/*
    SPDX-License-Identifier: GPL-3.0-or-later
    Copyright (c) 2025 WoozyMasta
    Upstream: https://github.com/WoozyMasta/metricz
    Modified: 2026
*/

#ifdef SERVER
/**
    \brief Utility for writing JSON files atomically.
    \details Writes JSON to a .tmp file, then copies it to the target .json file and removes the temp.
             If atomic mode is disabled in config, writes directly.
*/
class MetricZ_JsonWriter
{
	/**
	    \brief Write a JSON string atomically to exportDir/name.json.
	    \param exportDir Target directory (must end with /).
	    \param name File name without extension.
	    \param json JSON content to write.
	    \return true on success.
	*/
	static bool Write(string exportDir, string name, string json)
	{
		if (!MetricZ_Config.IsLoaded())
			return false;

		string target = exportDir + name + ".json";

		if (!MetricZ_Config.Get().json.atomic)
			return WriteFile(target, json);

		string tmp = exportDir + name + ".tmp";

		if (!WriteFile(tmp, json))
			return false;

		DeleteFile(target);

		if (!CopyFile(tmp, target)) {
			ErrorEx(
			    string.Format("MetricZ: atomic write failed '%1' -> '%2'", tmp, target),
			    ErrorExSeverity.ERROR);
			DeleteFile(tmp);
			return false;
		}

		DeleteFile(tmp);
		return true;
	}

	/**
	    \brief Write string to file, overwriting existing content.
	    \param path Full file path.
	    \param content Content to write.
	    \return true on success.
	*/
	protected static bool WriteFile(string path, string content)
	{
		FileHandle fh = OpenFile(path, FileMode.WRITE);
		if (!fh) {
			ErrorEx(
			    string.Format("MetricZ: open file '%1' failed", path),
			    ErrorExSeverity.ERROR);
			return false;
		}

		FPrint(fh, content);
		CloseFile(fh);
		return true;
	}

	/**
	    \brief Delete all .json and .tmp files in the export directory.
	    \param exportDir Target directory (must end with /).
	*/
	static void DeleteAll(string exportDir)
	{
		array<string> files = new array<string>();
		FindFileHandle fh = FindFile(exportDir + "*.json", files, FileAttr.NORMAL, FindFileFlags.CURRENT_DIR);
		if (fh) {
			foreach (string f : files)
				DeleteFile(exportDir + f);
			CloseFindFile(fh);
		}

		files.Clear();
		fh = FindFile(exportDir + "*.tmp", files, FileAttr.NORMAL, FindFileFlags.CURRENT_DIR);
		if (fh) {
			foreach (string f : files)
				DeleteFile(exportDir + f);
			CloseFindFile(fh);
		}
	}
}
#endif
